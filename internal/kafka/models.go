package kafka

import (
	"context"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
)

type Producer interface {
	SendMessage(ctx context.Context, order models.Order, event models.Event) error
	SendMessages(ctx context.Context, orders []models.Order, event models.Event) error
	Close() error
}

type Consumer interface {
	ConsumeMessage(ctx context.Context)
}

type Message struct {
	OrderID uint
	UserId  uint
	Event   models.Event
}
