package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opentracing/opentracing-go"
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
	parentSpan := opentracing.SpanFromContext(ctx)
	dbSpan := parentSpan.Tracer().StartSpan("AcceptOrder", opentracing.ChildOf(parentSpan.Context()))
	defer dbSpan.Finish()

	dbSpan.LogKV("action", "AcceptOrder", "order_id", or.ID)

	if or.KeepUntilDate.Before(time.Now()) {
		dbSpan.SetTag("error", true)
		dbSpan.LogKV("AcceptOrder", "error_message", e.ErrDate)
		return e.ErrDate
	}
	if _, t := s.cache.Get(ctx, strconv.Itoa(int(or.ID))); t {
		dbSpan.SetTag("error", true)
		dbSpan.LogKV("AcceptOrder", "error_message", e.ErrIsConsist)
		return e.ErrIsConsist
	}
	or.State = models.AcceptState
	or.AcceptTime = time.Now().Unix()
	err := s.PgRepo.AddToStorage(ctx, or)
	if err != nil {
		dbSpan.SetTag("error", true)
		dbSpan.LogKV("AcceptOrder", "error_message", err.Error())
		return err
	}

	s.cache.Set(ctx, strconv.Itoa(int(or.ID)), or, 5*time.Second)

	err = s.producer.SendMessage(ctx, *or, models.AcceptEvent)
	if err != nil {
		return ErrLogEvent
	}

	return nil
}

func (s storageFacade) PlaceOrder(ctx context.Context, ids []uint32) error {
	parentSpan := opentracing.SpanFromContext(ctx)
	span := parentSpan.Tracer().StartSpan("PlaceOrder", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

	span.LogKV("action", "PlaceOrder", "order_ids", ids)

	return s.txManager.RunReadCommited(ctx, func(ctxT context.Context) error {
		orders, exists := s.PgRepo.GetItems(ctxT, ids)
		if !exists {
			span.SetTag("error", true)
			span.LogKV("PlaceOrder", "error_message", e.ErrNoConsist)
			return e.ErrNoConsist
		}

		baseUserID := orders[0].UserID
		for _, order := range orders {
			if order.UserID != baseUserID {
				span.SetTag("error", true)
				span.LogKV("PlaceOrder", "error_message", e.ErrNotAllIDs)
				return e.ErrNotAllIDs
			}
			if order.State == models.PlaceState {
				err := fmt.Errorf("Order by id: %d is already placed", order.ID)
				span.SetTag("error", true)
				span.LogKV("PlaceOrder", "error_message", err.Error())
				return err
			}
			if order.State == models.SoftDelete {
				err := fmt.Errorf("Order by id: %d was deleted", order.ID)
				span.SetTag("error", true)
				span.LogKV("PlaceOrder", "error_message", err.Error())
				return err
			}
			if !order.KeepUntilDate.After(time.Now()) {
				err := fmt.Errorf("Order by id: %d cannot be issued to the customer because the date is invalid", order.ID)
				span.SetTag("error", true)
				span.LogKV("PlaceOrder", "error_message", err.Error())
				return err
			}
		}

		if err := s.PgRepo.UpdateBeforePlace(ctxT, ids, time.Now()); err != nil {
			span.SetTag("error", true)
			span.LogKV("PlaceOrder", "error_message", err.Error())
			return err
		}

		err := s.producer.SendMessages(ctx, orders, models.PlaceEvent)
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("PlaceOrder", "error_message", ErrLogEvent.Error())
			return ErrLogEvent
		}
		return nil
	})
}

func (s storageFacade) ReturnOrder(ctx context.Context, id uint) error {
	parentSpan := opentracing.SpanFromContext(ctx)
	span := parentSpan.Tracer().StartSpan("ReturnOrder", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()
	span.LogKV("action", "ReturnOrder", "order_id", id)

	return s.txManager.RunReadCommited(ctx, func(ctxTx context.Context) error {
		order, exists := s.PgRepo.GetItem(ctx, id)
		if !exists {
			span.SetTag("error", true)
			span.LogKV("ReturnOrder", "error_message", e.ErrNoConsist)
			return e.ErrNoConsist
		}
		err := order.CanBeReturned()
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("ReturnOrder", "error_message", err.Error())
			return err
		}
		err = s.PgRepo.DeleteFromStorage(ctxTx, order.ID)
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("ReturnOrder", "error_message", err.Error())
			return err
		}

		s.cache.Remove(strconv.Itoa(int(order.ID)))

		err = s.producer.SendMessage(ctx, *order, models.ReturnEvent)
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("ReturnOrder", "error_message", ErrLogEvent.Error())
			return ErrLogEvent
		}

		return nil
	})
}

func (s storageFacade) RefundOrder(ctx context.Context, id uint, userId uint) error {
	parentSpan := opentracing.SpanFromContext(ctx)
	span := parentSpan.Tracer().StartSpan("RefundOrder", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

	span.LogKV("action", "RefundOrder", "order_id", id, "user_id", userId)

	return s.txManager.RunReadCommited(ctx, func(ctxTx context.Context) error {
		order, exists := s.PgRepo.GetItem(ctx, id)
		if !exists {
			span.SetTag("error", true)
			span.LogKV("RefundOrder", "error_message", e.ErrCheckOrderID)
			return e.ErrCheckOrderID
		}
		if order.State != models.PlaceState {
			span.SetTag("error", true)
			span.LogKV("RefundOrder", "error_message", e.ErrNotPlace)
			return e.ErrNotPlace
		}
		if time.Now().After(order.PlaceDate.AddDate(0, 0, 2)) {
			span.SetTag("error", true)
			span.LogKV("RefundOrder", "error_message", e.ErrTimeExpired)
			return e.ErrTimeExpired
		}
		if order.UserID != userId {
			span.SetTag("error", true)
			span.LogKV("RefundOrder", "error_message", e.ErrIncorrectUserId)
			return e.ErrIncorrectUserId
		}

		err := s.PgRepo.UpdateState(ctxTx, order.ID, models.RefundedState)
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("RefundOrder", "error_message", err.Error())
			return err
		}

		err = s.producer.SendMessage(ctx, *order, models.RefundEvent)
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("RefundOrder", "error_message", ErrLogEvent.Error())
			return ErrLogEvent
		}

		return nil
	})
}

func (s storageFacade) ListOrders(ctx context.Context, id uint, inPuP bool, count int) ([]models.Order, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	span := parentSpan.Tracer().StartSpan("ListOrders", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

	span.LogKV("action", "ListOrders", "user_id", id, "in_progress", inPuP, "count", count)

	cacheKey := fmt.Sprintf("orders_%d_%t_%d", id, inPuP, count)
	if cachedOrders, found := s.cache.Get(ctx, cacheKey); found {
		return cachedOrders.([]models.Order), nil
	}

	list, err := s.PgRepo.GetUserOrders(ctx, id, inPuP)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("ListOrders", "error_message", err.Error())
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

	s.cache.Set(ctx, cacheKey, list, 5*time.Second)

	return list, nil
}

func (s storageFacade) ListReturns(ctx context.Context, limit, page int) ([]models.Order, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	span := parentSpan.Tracer().StartSpan("ListReturns", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

	span.LogKV("action", "ListReturns", "page", page, "limit", limit)

	cacheKey := fmt.Sprintf("returns_page_%d_limit_%d", page, limit)
	if cachedList, found := s.cache.Get(ctx, cacheKey); found {
		return cachedList.([]models.Order), nil
	}

	list, err := s.PgRepo.GetReturns(ctx, page-1, limit)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("ListReturns", "error_message", err.Error())
		return nil, err
	}

	s.cache.Set(ctx, cacheKey, list, 5*time.Second)

	return list, nil
}

func SortOrders(o []models.Order) {
	sort.Slice(o, func(i, j int) bool {
		return o[i].AcceptTime > o[j].AcceptTime
	})
}
