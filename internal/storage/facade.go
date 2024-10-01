package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
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
	GetItem(ctx context.Context, id uint) (*models.Order, bool)
}

type storageFacade struct {
	txManager postgres.TransactionManager
	PgRepo    *postgres.PgRepository
	PgReplica *postgres.PgRepository
}

func NewStorageFacade(pool *pgxpool.Pool) *storageFacade {
	txManager := postgres.NewTxManager(pool)
	PgRepo := postgres.NewPgRepository(txManager)

	return &storageFacade{
		txManager: txManager,
		PgRepo:    PgRepo,
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
	return nil
}

func (s storageFacade) PlaceOrder(ctx context.Context, ids []uint) error {
	return s.txManager.RunSerializable(ctx, func(ctxT context.Context) error {
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

		return nil
	})
}

func (s storageFacade) ReturnOrder(ctx context.Context, id uint) error {
	return s.txManager.RunSerializable(ctx, func(ctxTx context.Context) error {

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

		return nil
	})
}

func (s storageFacade) ListOrders(ctx context.Context, id uint, inPuP bool) ([]models.Order, error) {
	list, err := s.PgRepo.GetUserOrders(ctx, id, inPuP)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s storageFacade) RefundOrder(ctx context.Context, id uint, userId uint) error {
	return s.txManager.RunSerializable(ctx, func(ctxTx context.Context) error {
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
