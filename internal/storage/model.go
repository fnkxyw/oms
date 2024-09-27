package storage

import (
	"context"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"time"
)

type Storage interface {
	AddToStorage(ctx context.Context, order *models.Order)
	IsConsist(ctx context.Context, id uint) bool
	DeleteFromStorage(ctx context.Context, id uint)
	GetItem(ctx context.Context, id uint) (*models.Order, bool)
	GetIDs(ctx context.Context) []uint
	UpdateBeforePlace(ctx context.Context, id uint, state models.State, t time.Time) error
	UpdateState(ctx context.Context, id uint, state models.State) error
}
