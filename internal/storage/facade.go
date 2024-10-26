package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/cache"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/kafka"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/postgres"
	"sort"
	"strconv"
	"time"
)

var (
	ErrLogEvent = fmt.Errorf("error logging event")
)

type Facade interface {
	AcceptOrder(ctx context.Context, or *models.Order) error
	PlaceOrder(ctx context.Context, ids []uint32) error
	ReturnOrder(ctx context.Context, id uint) error
	ListOrders(ctx context.Context, id uint, inPuP bool, count int) ([]models.Order, error)
	RefundOrder(ctx context.Context, id uint, userId uint) error
	ListReturns(ctx context.Context, limit, page int) ([]models.Order, error)
}

type storageFacade struct {
	producer  kafka.Producer
	txManager postgres.TransactionManager
	PgRepo    *postgres.PgRepository
	cache     *cache.Cache[string, any]
}

func NewStorageFacade(pool *pgxpool.Pool, producer kafka.Producer) *storageFacade {
	txManager := postgres.NewTxManager(pool)
	PgRepo := postgres.NewPgRepository(txManager)
	ch := cache.NewCache[string, any](5)
	return &storageFacade{
		txManager: txManager,
		PgRepo:    PgRepo,
		producer:  producer,
		cache:     ch,
	}
}

func (s storageFacade) AcceptOrder(ctx context.Context, or *models.Order) error {
	if or.KeepUntilDate.Before(time.Now()) {
		return e.ErrDate
	}

	if _, t := s.cache.Get(strconv.Itoa(int(or.ID))); t {
		return e.ErrIsConsist
	}
	or.State = models.AcceptState
	or.AcceptTime = time.Now().Unix()
	err := s.PgRepo.AddToStorage(ctx, or)
	if err != nil {
		return err
	}

	s.cache.Set(strconv.Itoa(int(or.ID)), or, 5*time.Second)

	err = s.producer.SendMessage(ctx, *or, models.AcceptEvent)
	if err != nil {
		return ErrLogEvent
	}

	return nil
}

func (s storageFacade) PlaceOrder(ctx context.Context, ids []uint32) error {
	return s.txManager.RunReadCommited(ctx, func(ctxT context.Context) error {
		var orders []models.Order
		var missingIDs []uint32
		for _, id := range ids {
			if order, found := s.cache.Get(strconv.Itoa(int(id))); found {
				orders = append(orders, *order.(*models.Order))
			} else {
				missingIDs = append(missingIDs, id)
			}
		}

		if len(missingIDs) > 0 {
			dbOrders, exists := s.PgRepo.GetItems(ctxT, missingIDs)
			if !exists {
				return e.ErrNoConsist
			}
			for _, order := range dbOrders {
				s.cache.Set(strconv.Itoa(int(order.ID)), &order, 30*time.Second)
				orders = append(orders, order)
			}
		}
		baseUserID := orders[0].UserID
		for _, order := range orders {
			if order.UserID != baseUserID {
				return e.ErrNotAllIDs
			}

			switch order.State {
			case models.PlaceState:
				return fmt.Errorf("Order by id: %d is already placed", order.ID)
			case models.SoftDelete:
				return fmt.Errorf("Order by id: %d was deleted", order.ID)
			}

			if !order.KeepUntilDate.After(time.Now()) {
				return fmt.Errorf("Order by id: %d cannot be issued to the customer because the date is invalid", order.ID)
			}
		}

		if err := s.PgRepo.UpdateBeforePlace(ctxT, ids, time.Now()); err != nil {
			return err
		}
		err := s.producer.SendMessages(ctx, orders, models.PlaceEvent)
		if err != nil {
			return ErrLogEvent
		}
		return nil
	})
}

func (s storageFacade) ReturnOrder(ctx context.Context, id uint) error {
	return s.txManager.RunReadCommited(ctx, func(ctxTx context.Context) error {
		var order *models.Order
		if cachedOrder, t := s.cache.Get(strconv.Itoa(int(id))); t {
			order = cachedOrder.(*models.Order)
		} else {
			var exists bool
			order, exists = s.PgRepo.GetItem(ctxTx, id)
			if !exists {
				return e.ErrNoConsist
			}
			s.cache.Set(strconv.Itoa(int(order.ID)), order, 5*time.Second)
		}

		err := order.CanBeReturned()
		if err != nil {
			return err
		}
		err = s.PgRepo.DeleteFromStorage(ctxTx, order.ID)
		if err != nil {
			return err
		}

		s.cache.Remove(strconv.Itoa(int(order.ID)))

		err = s.producer.SendMessage(ctx, *order, models.ReturnEvent)
		if err != nil {
			return ErrLogEvent
		}

		return nil
	})
}

func (s storageFacade) ListOrders(ctx context.Context, id uint, inPuP bool, count int) ([]models.Order, error) {
	cacheKey := fmt.Sprintf("orders_%d_%t_%d", id, inPuP, count)

	if cachedOrders, found := s.cache.Get(cacheKey); found {
		return cachedOrders.([]models.Order), nil
	}

	list, err := s.PgRepo.GetUserOrders(ctx, id, inPuP)
	if err != nil {
		return nil, err
	}
	SortOrders(list)

	if count < 1 {
		count = 1
	} else if count > len(list) {
		count = len(list)
	}

	if !inPuP {
		list = list[:count]
	}

	s.cache.Set(cacheKey, list, 5*time.Second)

	return list, nil
}

func (s storageFacade) RefundOrder(ctx context.Context, id uint, userId uint) error {
	return s.txManager.RunReadCommited(ctx, func(ctxTx context.Context) error {
		var order *models.Order
		var exists bool

		if cachedOrder, found := s.cache.Get(strconv.Itoa(int(id))); found {
			order = cachedOrder.(*models.Order)
		} else {
			order, exists = s.PgRepo.GetItem(ctxTx, id)
			if !exists {
				return e.ErrCheckOrderID
			}
		}

		if order.State != models.PlaceState {
			return e.ErrNotPlace
		}
		if time.Now().After(order.PlaceDate.AddDate(0, 0, 2)) {
			return e.ErrTimeExpired
		}
		if order.UserID != userId {
			return e.ErrIncorrectUserId
		}

		err := s.PgRepo.UpdateState(ctxTx, id, models.RefundedState)
		if err != nil {
			return err
		}

		s.cache.Set(strconv.Itoa(int(order.ID)), order, 5*time.Second)

		err = s.producer.SendMessage(ctx, *order, models.RefundEvent)
		if err != nil {
			return ErrLogEvent
		}

		return nil
	})
}

func (s storageFacade) ListReturns(ctx context.Context, limit, page int) ([]models.Order, error) {
	cacheKey := fmt.Sprintf("returns_page_%d_limit_%d", page, limit)

	if cachedList, found := s.cache.Get(cacheKey); found {
		return cachedList.([]models.Order), nil
	}

	list, err := s.PgRepo.GetReturns(ctx, page-1, limit)
	if err != nil {
		return nil, err
	}

	s.cache.Set(cacheKey, list, 5*time.Second)

	return list, nil
}

func SortOrders(o []models.Order) {
	sort.Slice(o, func(i, j int) bool {
		return o[i].AcceptTime > o[j].AcceptTime
	})
}
