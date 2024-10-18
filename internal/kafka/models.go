package kafka

import "gitlab.ozon.dev/akugnerevich/homework.git/internal/models"

type KafkaProducer interface {
	SendMessage(order models.Order, state models.Event) error
	SendMessages(orders []models.Order, event models.Event) error
}

type Message struct {
	OrderID uint
	UserId  uint
	Event   models.Event
}
