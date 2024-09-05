package storage

import (
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
)

type ReturnStorage struct {
	Data map[uint]*models.Return
}

func NewReturnStorage() *ReturnStorage {
	return &ReturnStorage{Data: make(map[uint]*models.Return)}
}

func (rs *ReturnStorage) AddReturnToStorage(r *models.Return) error {
	_, ok := rs.Data[r.ID]
	if ok {
		return fmt.Errorf("Order already return ")
	} else {
		rs.Data[r.ID] = r
	}

	return nil
}

func (r *ReturnStorage) DeleteOrderFromStorage(id uint) {
	delete(r.Data, id)
}
