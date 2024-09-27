package storage

import (
	"context"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"time"
)

type Storage interface {
	AddToStorage(ctx context.Context, order *models.Order)
	DeleteFromStorage(ctx context.Context, id uint)

	IsConsist(ctx context.Context, id uint) bool

	GetItem(ctx context.Context, id uint) (*models.Order, bool)
	GetIDs(ctx context.Context) []uint
	GetByUserId(ctx context.Context, userId uint) ([]*models.Order, error)
	GetReturns(ctx context.Context, state models.State) ([]*models.Order, error)

	UpdateBeforePlace(ctx context.Context, id uint, state models.State, t time.Time) error
	UpdateState(ctx context.Context, id uint, state models.State) error
}
