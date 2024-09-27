package orderStorage

import (
	"context"
	"encoding/json"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/inmemory/errors"
	"io"
	"os"
	"time"
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

func (o *OrderStorage) AddToStorage(ctx context.Context, order *models.Order) {
	o.Data[order.ID] = order

}

func (o *OrderStorage) IsConsist(ctx context.Context, id uint) bool {
	_, ok := o.Data[id]
	return ok
}

func (o *OrderStorage) DeleteFromStorage(ctx context.Context, id uint) {
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

func (o *OrderStorage) GetItem(ctx context.Context, id uint) (*models.Order, bool) {
	temp, ok := o.Data[id]
	return temp, ok
}

func (o *OrderStorage) GetIDs(ctx context.Context) []uint {
	var ids []uint
	for id := range o.Data {
		ids = append(ids, id)
	}
	return ids
}

func (o *OrderStorage) GetPath(ctx context.Context) string {
	return o.path
}

func (o *OrderStorage) SetPath(ctx context.Context, p string) {
	o.path = p
}

//UpdateBeforePlace(ctx context.Context, id uint, state models.State, t time.Time) error
//UpdateState(ctx context.Context, id uint, state models.State) error

func (o *OrderStorage) UpdateBeforePlace(ctx context.Context, id uint, state models.State, t time.Time) error {
	order := o.Data[id]
	order.State = state
	order.PlaceDate = t
	return nil
}

func (o *OrderStorage) UpdateState(ctx context.Context, id uint, state models.State) error {
	order := o.Data[id]
	order.State = state
	return nil
}

func (o *OrderStorage) GetByUserId(ctx context.Context, userId uint) ([]*models.Order, error) {
	var orders []*models.Order
	for _, order := range o.Data {
		if order.UserID == userId {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (o *OrderStorage) GetReturns(ctx context.Context, state models.State) ([]*models.Order, error) {
	var orders []*models.Order
	for _, order := range o.Data {
		if order.State == state {
			orders = append(orders, order)
		}
	}
	return orders, nil
}
