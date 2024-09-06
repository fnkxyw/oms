package storage

import (
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
	"io"
	"os"
)

type OrderStorage struct {
	Data map[uint]*models.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{Data: make(map[uint]*models.Order)}
}

func (os *OrderStorage) AddOrderToStorage(or *models.Order) error {
	_, ok := os.Data[or.ID]
	if ok {
		return fmt.Errorf("Order already accept\n")
	} else {
		os.Data[or.ID] = or
	}

	return nil
}

func (o *OrderStorage) IsConsist(id uint) bool {
	_, ok := o.Data[id]
	if !ok {
		return false
	}
	return true
}

func (o *OrderStorage) DeleteOrderFromStorage(id uint) {
	delete(o.Data, id)
}

// считываем с JSON-a
func (o *OrderStorage) ReadFromJSON(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("Open file error: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Read file error: %w", err)
	}

	if len(data) == 0 {
		o.Data = make(map[uint]*models.Order)
		return nil
	}

	var i struct {
		Data map[uint]models.Order `json:"Data"`
	}

	err = json.Unmarshal(data, &i)
	if err != nil {
		return fmt.Errorf("ошибка при декодировании JSON: %w", err)
	}

	o.Data = make(map[uint]*models.Order)
	for orderID, order := range i.Data {
		orderCopy := order
		o.Data[orderID] = &orderCopy
	}

	return nil
}
