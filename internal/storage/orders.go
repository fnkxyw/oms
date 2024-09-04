package storage

import (
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
)

type OrderStorage struct {
	Data map[uint]map[uint]*models.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{Data: make(map[uint]map[uint]*models.Order)}
}

func (os *OrderStorage) AddOrderToStorage(or *models.Order) error {
	_, ok := os.Data[or.ID][or.UserID]
	if ok {
		return fmt.Errorf("Order already accept")
	} else {
		os.Data[or.ID] = make(map[uint]*models.Order)
	}

	os.Data[or.ID][or.UserID] = or
	return nil
}
