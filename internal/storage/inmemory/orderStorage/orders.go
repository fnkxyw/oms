package orderStorage

import (
	"encoding/json"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/inmemory/errors"
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
	return &OrderStorage{Data: make(map[uint]*models.Order), path: "api/order.json"}
}

func (o *OrderStorage) AddToStorage(order *models.Order) {
	o.Data[order.ID] = order

}

func (o *OrderStorage) IsConsist(id uint) bool {
	_, ok := o.Data[id]
	return ok
}

func (o *OrderStorage) DeleteFromStorage(id uint) {
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

func (o *OrderStorage) GetItem(id uint) (*models.Order, bool) {
	temp, ok := o.Data[id]
	return temp, ok
}

func (o *OrderStorage) GetIDs() []uint {
	var ids []uint
	for id := range o.Data {
		ids = append(ids, id)
	}
	return ids
}

func (o *OrderStorage) GetPath() string {
	return o.path
}

func (o *OrderStorage) SetPath(p string) {
	o.path = p
}
