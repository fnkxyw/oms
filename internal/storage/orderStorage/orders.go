package orderStorage

import (
	"encoding/json"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/errors"
	"io"
	"os"
)

type OrderStorageInterface interface {
	AddOrderToStorage(or *models.Order)
	IsConsist(id uint) bool
	DeleteOrderFromStorage(id uint)
	GetOrder(id uint) (*models.Order, bool)
	GetOrderIDs() []uint
	ReadFromJSON() error
	WriteToJSON() error
}

type OrderStorage struct {
	Data map[uint]*models.Order
	path string
}

func (o *OrderStorage) Create() error {
	_, err := os.Create(o.path)
	return err
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{Data: make(map[uint]*models.Order), path: "api/order.json"}
}

func (os *OrderStorage) AddOrderToStorage(or *models.Order) {
	os.Data[or.ID] = or
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
	file, err := os.OpenFile(o.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return e.ErrOpenFile
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return e.ErrReadFile
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
		return err
	}

	o.Data = make(map[uint]*models.Order)
	for orderID, order := range i.Data {
		orderCopy := order
		o.Data[orderID] = &orderCopy
	}

	return nil
}

func (o *OrderStorage) WriteToJSON() error {
	file, err := os.OpenFile(o.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return e.ErrOpenFile
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent(" ", "  ")
	if err := encoder.Encode(o); err != nil {
		return e.ErrEnocde
	}
	return nil
}

func (os *OrderStorage) GetOrder(id uint) (*models.Order, bool) {
	order, ok := os.Data[id]
	return order, ok
}

func (os *OrderStorage) GetOrderIDs() []uint {
	var ids []uint
	for id := range os.Data {
		ids = append(ids, id)
	}
	return ids
}

func (os *OrderStorage) GetPath() string {
	return os.path
}

func (os *OrderStorage) SetPath(p string) {
	os.path = p
}
