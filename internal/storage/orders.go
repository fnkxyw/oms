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
	path string
}

func (o *OrderStorage) Create() error {
	_, err := os.Create(o.path)
	return err
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{Data: make(map[uint]*models.Order), path: "api/orders.json"}
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
	return ok
}

func (o *OrderStorage) DeleteOrderFromStorage(id uint) {
	delete(o.Data, id)
}

// считываем с JSON-a
func (o *OrderStorage) ReadFromJSON() error {
	file, err := os.OpenFile(o.path, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("Open file Err: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Read file Err: %w", err)
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

func (o *OrderStorage) WritoToJSON() error {
	file, err := os.OpenFile("api/orders.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("OpenFile eror in WriteToJSON", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent(" ", "  ")
	if err := encoder.Encode(o); err != nil {
		fmt.Println("Encoding Err in WirteToJSON", err)
		return err
	}
	return nil
}
