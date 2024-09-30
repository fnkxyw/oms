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
	ListOrders(ctx context.Context, id uint, inPuP bool) ([]models.Order, error)
	RefundOrder(ctx context.Context, id uint, userId uint) error
	ListReturns(ctx context.Context, limit, page int) ([]models.Order, error)
	CheckIDsOrders(ctx context.Context, ids []uint) error
}

type storageFacade struct {
	txManager    postgres.TransactionManager
	pgRepository *postgres.PgRepository
	pgReplica    *postgres.PgRepository
}

func NewStorageFacade(
	txManager postgres.TransactionManager,
	pgRepository *postgres.PgRepository,
	pgReplica *postgres.PgRepository,
) *storageFacade {
	return &storageFacade{
		txManager:    txManager,
		pgRepository: pgRepository,
		pgReplica:    pgReplica,
	}
}

func (s storageFacade) AcceptOrder(ctx context.Context, or *models.Order) error {
	return s.txManager.RunReadCommited(ctx, func(ctxTx context.Context) error {

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
	})
}

func (s storageFacade) PlaceOrder(ctx context.Context, ids []uint) error {
	return s.txManager.RunSerializable(ctx, func(ctxT context.Context) error {
		orders, exists := s.pgRepository.GetItems(ctxT, ids)
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

		if err := s.pgRepository.UpdateBeforePlace(ctxT, ids, time.Now()); err != nil {
			return err
		}

		return nil
	})
}

func (s storageFacade) ReturnOrder(ctx context.Context, id uint) error {
	return s.txManager.RunSerializable(ctx, func(ctxTx context.Context) error {

		order, exists := s.pgReplica.GetItem(ctx, id)
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

func (s storageFacade) ListOrders(ctx context.Context, id uint, inPuP bool) ([]models.Order, error) {
	var list []models.Order
	list, err := s.pgRepository.GetUserOrders(ctx, id, inPuP)
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
