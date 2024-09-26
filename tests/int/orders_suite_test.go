package int

import (
	"encoding/json"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/orderStorage"
)

type OrderSuite struct {
	suite.Suite
	storage *orderStorage.OrderStorage
}

func (s *OrderSuite) SetupTest() {
	s.storage = orderStorage.NewOrderStorage()
	s.storage.SetPath("order_test.json")
}

func (s *OrderSuite) TestAcceptOrder() {
	order := &models.Order{
		ID:            1,
		UserID:        1,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}

	err := orders.AcceptOrder(s.storage, order)
	s.NoError(err)

	storedOrder, exists := s.storage.GetOrder(order.ID)
	s.True(exists)
	s.Equal(models.AcceptState, storedOrder.State)
	s.NotZero(storedOrder.AcceptTime)
}

func (s *OrderSuite) TestAcceptOrder_PastDate() {
	order := &models.Order{
		ID:            3,
		UserID:        1,
		KeepUntilDate: time.Now().Add(-1 * time.Hour),
	}

	err := orders.AcceptOrder(s.storage, order)
	s.ErrorIs(err, e.ErrDate)
}

func (s *OrderSuite) TestAcceptOrder_EqualOrder() {
	order1 := &models.Order{
		ID:            1,
		UserID:        1,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}
	order2 := &models.Order{
		ID:            1,
		UserID:        1,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}
	err := orders.AcceptOrder(s.storage, order1)
	s.ErrorIs(err, nil)
	err = orders.AcceptOrder(s.storage, order2)
	s.ErrorIs(err, e.ErrIsConsist)
}

func (s *OrderSuite) TestPlaceOrder() {
	order := &models.Order{
		ID:            2,
		UserID:        1,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
		State:         models.AcceptState,
	}
	s.storage.AddOrderToStorage(order)

	err := orders.PlaceOrder(s.storage, []uint{2})
	s.NoError(err)

	updatedOrder, exists := s.storage.GetOrder(2)
	s.True(exists)
	s.Equal(models.PlaceState, updatedOrder.State)
}

func (s *OrderSuite) TestPlaceOrder_NoCount() {
	err := orders.PlaceOrder(s.storage, []uint{1})
	s.ErrorIs(err, e.ErrNoConsist)
}

func (s *OrderSuite) TestListOrder() {
	order1 := &models.Order{
		ID:            4,
		UserID:        1,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}
	order2 := &models.Order{
		ID:            5,
		UserID:        1,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}
	s.storage.AddOrderToStorage(order1)
	s.storage.AddOrderToStorage(order2)

	err := orders.ListOrders(s.storage, 1, 2, false)
	s.NoError(err)

	ids := s.storage.GetOrderIDs()
	s.Len(ids, 2)
	s.Contains(ids, uint(4))
	s.Contains(ids, uint(5))
}

func (s *OrderSuite) TestReturnOrder() {
	order := &models.Order{
		ID:     6,
		UserID: 1,
		State:  models.ReturnedState,
	}
	s.storage.AddOrderToStorage(order)

	err := orders.ReturnOrder(s.storage, 6)
	s.NoError(err)

	s.Error(order.CanReturned())
}

func (s *OrderSuite) TestReturnOrder_NoCount() {
	order := &models.Order{
		ID:     1,
		UserID: 1,
		State:  models.ReturnedState,
	}
	err := orders.ReturnOrder(s.storage, order.ID)
	s.ErrorIs(err, e.ErrNoConsist)
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}

func (s *OrderSuite) TestFilterOrder() {
	order1 := &models.Order{
		ID:            7,
		UserID:        1,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}
	order2 := &models.Order{
		ID:            8,
		UserID:        1,
		State:         models.ReturnedState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}
	order3 := &models.Order{
		ID:            9,
		UserID:        2,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}
	s.storage.AddOrderToStorage(order1)
	s.storage.AddOrderToStorage(order2)
	s.storage.AddOrderToStorage(order3)

	filteredOrder := orders.FilterOrders(s.storage, 1, true)
	s.Len(filteredOrder, 2)

	s.Contains(filteredOrder, order1)
	s.Contains(filteredOrder, order2)
}

func (s *OrderSuite) TestCheckIDOrder() {
	order1 := &models.Order{
		ID:     10,
		UserID: 1,
		State:  models.AcceptState,
	}
	order2 := &models.Order{
		ID:     11,
		UserID: 1,
		State:  models.AcceptState,
	}
	order3 := &models.Order{
		ID:     12,
		UserID: 2,
		State:  models.AcceptState,
	}
	s.storage.AddOrderToStorage(order1)
	s.storage.AddOrderToStorage(order2)
	s.storage.AddOrderToStorage(order3)

	err := orders.CheckIDsOrders(s.storage, []uint{10, 11})
	s.NoError(err)

	err = orders.CheckIDsOrders(s.storage, []uint{10, 12})
	s.ErrorIs(err, e.ErrNotAllIDs)
}

func (s *OrderSuite) TestWriteToJSON() {
	order := &models.Order{
		ID:            1,
		UserID:        1,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}

	s.storage.AddOrderToStorage(order)

	err := s.storage.WriteToJSON()
	s.NoError(err)

	file, err := os.Open(s.storage.GetPath())
	s.NoError(err)
	defer file.Close()

	os.Remove(s.storage.GetPath())
	var storageData orderStorage.OrderStorage
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&storageData)
	s.NoError(err)

	storedOrder, exists := storageData.GetOrder(order.ID)
	s.True(exists)
	s.Equal(order.UserID, storedOrder.UserID)
	s.Equal(order.State, storedOrder.State)
	s.Equal(order.KeepUntilDate.Unix(), storedOrder.KeepUntilDate.Unix())
}

func (s *OrderSuite) TestReadFromJSON() {
	order := &models.Order{
		ID:            1,
		UserID:        1,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}

	s.storage.AddOrderToStorage(order)

	err := s.storage.WriteToJSON()
	s.NoError(err)

	s.storage.Data = make(map[uint]*models.Order)

	err = s.storage.ReadFromJSON()
	s.NoError(err)

	os.Remove(s.storage.GetPath())

	storedOrder, exists := s.storage.GetOrder(order.ID)
	s.True(exists)
	s.Equal(order.UserID, storedOrder.UserID)
	s.Equal(order.State, storedOrder.State)
	s.Equal(order.KeepUntilDate.Unix(), storedOrder.KeepUntilDate.Unix())
}
