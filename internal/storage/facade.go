package storage

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/postgres"
	"time"
)

type Facade interface {
	AcceptOrder(ctx context.Context, or *models.Order) error
	PlaceOrder(ctx context.Context, ids []uint) error
	ReturnOrder(ctx context.Context, id uint) error
	ListOrders(ctx context.Context, id uint, n int, inPuP bool) ([]models.Order, error)
	RefundOrder(ctx context.Context, id uint, userId uint) error
	ListReturns(ctx context.Context, limit, page int) ([]models.Order, error)
	CheckIDsOrders(ctx context.Context, ids []uint) error
}

type storageFacade struct {
	txManager    postgres.TransactionManager
	pgRepository *postgres.PgRepository
}

func NewStorageFacade(
	txManager postgres.TransactionManager,
	pgRepository *postgres.PgRepository,
) *storageFacade {
	return &storageFacade{
		txManager:    txManager,
		pgRepository: pgRepository,
	}
}

func (s storageFacade) AcceptOrder(ctx context.Context, or *models.Order) error {
	if s.pgRepository.IsConsist(ctx, or.ID) {
		return e.ErrIsConsist
	}
	if or.KeepUntilDate.Before(time.Now()) {
		return e.ErrDate
	}
	or.State = models.AcceptState
	or.AcceptTime = time.Now().Unix()
	err := s.pgRepository.AddToStorage(ctx, or)
	if err != nil {
		return err
	}
	return nil
}

func (s storageFacade) PlaceOrder(ctx context.Context, ids []uint) error {
	return s.txManager.RunSerializable(ctx, func(ctxT context.Context) error {
		for _, id := range ids {
			order, exists := s.pgRepository.GetItem(ctx, id)
			if !exists {
				return e.ErrNoConsist
			}
			if order.State == models.PlaceState {
				return fmt.Errorf("Order by id: %d is already placed", id)
			}

			if order.State == models.SoftDelete {
				return fmt.Errorf("Order by id: %d was deleted", id)
			}

			if !order.KeepUntilDate.After(time.Now()) {
				return fmt.Errorf("Order by id: %d cannot be issued to the customer because the date is invalid", id)
			}

			err := s.pgRepository.UpdateBeforePlace(ctx, id, models.PlaceState, time.Now())
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s storageFacade) ReturnOrder(ctx context.Context, id uint) error {
	return s.txManager.RunSerializable(ctx, func(ctxTx context.Context) error {

		order, exists := s.pgRepository.GetItem(ctx, id)
		if !exists {
			return e.ErrNoConsist
		}
		err := order.CanReturned()
		if err != nil {
			return err
		}
		err = s.pgRepository.DeleteFromStorage(ctx, order.ID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s storageFacade) ListOrders(ctx context.Context, id uint, n int, inPuP bool) ([]models.Order, error) {
	var list []models.Order
	list, err := s.pgRepository.GetOrders(ctx, id, inPuP)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s storageFacade) RefundOrder(ctx context.Context, id uint, userId uint) error {
	return s.txManager.RunSerializable(ctx, func(ctxTx context.Context) error {
		order, exists := s.pgRepository.GetItem(ctx, id)

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

		err := s.pgRepository.UpdateState(ctx, id, models.RefundedState)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s storageFacade) ListReturns(ctx context.Context, limit, page int) ([]models.Order, error) {
	var list []models.Order
	list, err := s.pgRepository.GetReturns(ctx, page-1, limit)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		fmt.Printf("OrderID: %d, UserID: %d \n", v.ID, v.UserID)
	}
	return list, nil
}

func (s storageFacade) CheckIDsOrders(ctx context.Context, ids []uint) error {
	order, ok := s.pgRepository.GetItem(ctx, ids[0])
	if !ok {
		return e.ErrNoConsist
	}
	temp := order.UserID
	for _, id := range ids {
		order, _ = s.pgRepository.GetItem(ctx, id)
		if order.UserID != temp {
			return e.ErrNotAllIDs
		}
	}
	return nil
}
