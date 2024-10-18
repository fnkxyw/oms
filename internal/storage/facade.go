package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/kafka"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/postgres"
	"sort"
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
	GetItem(ctx context.Context, id uint) (*models.Order, bool)
}

type storageFacade struct {
	producer  kafka.KafkaProducer
	txManager postgres.TransactionManager
	PgRepo    *postgres.PgRepository
	PgReplica *postgres.PgRepository
}

func NewStorageFacade(pool *pgxpool.Pool, producer kafka.KafkaProducer) *storageFacade {
	txManager := postgres.NewTxManager(pool)
	PgRepo := postgres.NewPgRepository(txManager)

	return &storageFacade{
		txManager: txManager,
		PgRepo:    PgRepo,
		producer:  producer,
	}
}

func (s storageFacade) AcceptOrder(ctx context.Context, or *models.Order) error {
	if or.KeepUntilDate.Before(time.Now()) {
		return e.ErrDate
	}
	or.State = models.AcceptState
	or.AcceptTime = time.Now().Unix()
	err := s.PgRepo.AddToStorage(ctx, or)
	if err != nil {
		return err
	}

	err = s.producer.SendMessage(*or, models.AcceptEvent)
	if err != nil {
		return ErrLogEvent
	}

	return nil
}

func (s storageFacade) PlaceOrder(ctx context.Context, ids []uint32) error {
	return s.txManager.RunReadCommited(ctx, func(ctxT context.Context) error {
		orders, exists := s.PgRepo.GetItems(ctxT, ids)
		if !exists {
			return e.ErrNoConsist
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
		err := s.producer.SendMessages(orders, models.PlaceEvent)
		if err != nil {
			return ErrLogEvent
		}
		return nil
	})
}

func (s storageFacade) ReturnOrder(ctx context.Context, id uint) error {
	return s.txManager.RunReadCommited(ctx, func(ctxTx context.Context) error {

		order, exists := s.PgRepo.GetItem(ctx, id)
		if !exists {
			return e.ErrNoConsist
		}
		err := order.CanBeReturned()
		if err != nil {
			return err
		}
		err = s.PgRepo.DeleteFromStorage(ctx, order.ID)
		if err != nil {
			return err
		}
		err = s.producer.SendMessage(*order, models.ReturnEvent)
		if err != nil {
			return ErrLogEvent
		}

		return nil
	})
}

func (s storageFacade) ListOrders(ctx context.Context, id uint, inPuP bool, count int) ([]models.Order, error) {
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

	return list, nil
}

func (s storageFacade) RefundOrder(ctx context.Context, id uint, userId uint) error {
	return s.txManager.RunReadCommited(ctx, func(ctxTx context.Context) error {
		order, exists := s.PgRepo.GetItem(ctx, id)

		if !exists {
			return e.ErrCheckOrderID
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

		err := s.PgRepo.UpdateState(ctx, id, models.RefundedState)
		if err != nil {
			return err
		}
		err = s.producer.SendMessage(*order, models.RefundEvent)
		if err != nil {
			return ErrLogEvent
		}

		return nil
	})
}

func (s storageFacade) ListReturns(ctx context.Context, limit, page int) ([]models.Order, error) {
	var list []models.Order
	list, err := s.PgRepo.GetReturns(ctx, page-1, limit)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s storageFacade) GetItem(ctx context.Context, id uint) (*models.Order, bool) {
	return s.PgRepo.GetItem(ctx, id)
}

func SortOrders(o []models.Order) {
	sort.Slice(o, func(i, j int) bool {
		return o[i].AcceptTime > o[j].AcceptTime
	})
}
