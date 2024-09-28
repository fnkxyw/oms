package int

import (
	"context"
	"encoding/json"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	or "gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/inmemory/orderStorage"
)

type OrderSuite struct {
	suite.Suite
	storage *or.OrderStorage
}

func (s *OrderSuite) SetupTest() {
	ctx := context.Background()

	s.storage = or.NewOrderStorage()
	s.storage.SetPath(ctx, "order_test.json")
}

func (s *OrderSuite) TestAcceptOrder() {
	ctx := context.Background()

	order := &models.Order{
		ID:            1,
		UserID:        1,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}

	err := orders.AcceptOrder(ctx, s.storage, order)
	s.NoError(err)

	storedOrder, exists := s.storage.GetItem(ctx, order.ID)
	s.True(exists)
	s.Equal(models.AcceptState, storedOrder.State)
	s.NotZero(storedOrder.AcceptTime)
}

func (s *OrderSuite) TestAcceptOrder_PastDate() {
	ctx := context.Background()

	order := &models.Order{
		ID:            3,
		UserID:        1,
		KeepUntilDate: time.Now().Add(-1 * time.Hour),
	}

	err := orders.AcceptOrder(ctx, s.storage, order)
	s.ErrorIs(err, e.ErrDate)
}

func (s *OrderSuite) TestAcceptOrder_EqualOrder() {
	ctx := context.Background()

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
	err := orders.AcceptOrder(ctx, s.storage, order1)
	s.ErrorIs(err, nil)
	err = orders.AcceptOrder(ctx, s.storage, order2)
	s.ErrorIs(err, e.ErrIsConsist)
}

func (s *OrderSuite) TestPlaceOrder() {
	ctx := context.Background()

	order := &models.Order{
		ID:            2,
		UserID:        1,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
		State:         models.AcceptState,
	}
	s.storage.AddToStorage(ctx, order)

	err := orders.PlaceOrder(ctx, s.storage, []uint{2})
	s.NoError(err)

	updatedOrder, exists := s.storage.GetItem(ctx, 2)
	s.True(exists)
	s.Equal(models.PlaceState, updatedOrder.State)
}

func (s *OrderSuite) TestPlaceOrder_NoCount() {
	ctx := context.Background()

	err := orders.PlaceOrder(ctx, s.storage, []uint{1})
	s.ErrorIs(err, e.ErrNoConsist)
}

func (s *OrderSuite) TestListOrder() {
	ctx := context.Background()

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
	s.storage.AddToStorage(ctx, order1)
	s.storage.AddToStorage(ctx, order2)

	err := orders.ListOrders(ctx, s.storage, 1, 2, false)
	s.NoError(err)

	ids := s.storage.GetIDs(ctx)
	s.Len(ids, 2)
	s.Contains(ids, uint(4))
	s.Contains(ids, uint(5))
}

func (s *OrderSuite) TestReturnOrder() {
	ctx := context.Background()

	order := &models.Order{
		ID:     6,
		UserID: 1,
		State:  models.RefundedState,
	}
	s.storage.AddToStorage(ctx, order)

	err := orders.ReturnOrder(ctx, s.storage, 6)
	s.NoError(err)

	s.NoError(order.CanReturned())
}

func (s *OrderSuite) TestReturnOrder_NoCount() {
	ctx := context.Background()

	order := &models.Order{
		ID:     1,
		UserID: 1,
		State:  models.RefundedState,
	}
	err := orders.ReturnOrder(ctx, s.storage, order.ID)
	s.ErrorIs(err, e.ErrNoConsist)
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}

func (s *OrderSuite) TestFilterOrder() {
	ctx := context.Background()

	order1 := &models.Order{
		ID:            7,
		UserID:        1,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}
	order2 := &models.Order{
		ID:            8,
		UserID:        1,
		State:         models.RefundedState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}
	order3 := &models.Order{
		ID:            9,
		UserID:        2,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}
	s.storage.AddToStorage(ctx, order1)
	s.storage.AddToStorage(ctx, order2)
	s.storage.AddToStorage(ctx, order3)

	filteredOrder := orders.FilterOrders(ctx, s.storage, 1, true)
	s.Len(filteredOrder, 2)

	s.Contains(filteredOrder, order1)
	s.Contains(filteredOrder, order2)
}

func (s *OrderSuite) TestCheckIDOrder() {
	ctx := context.Background()

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
	s.storage.AddToStorage(ctx, order1)
	s.storage.AddToStorage(ctx, order2)
	s.storage.AddToStorage(ctx, order3)

	err := orders.CheckIDsOrders(ctx, s.storage, []uint{10, 11})
	s.NoError(err)

	err = orders.CheckIDsOrders(ctx, s.storage, []uint{10, 12})
	s.ErrorIs(err, e.ErrNotAllIDs)
}

func (s *OrderSuite) TestWriteToJSON() {
	ctx := context.Background()

	order := &models.Order{
		ID:            1,
		UserID:        1,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}

	s.storage.AddToStorage(ctx, order)

	err := s.storage.WriteToJSON()
	s.NoError(err)

	file, err := os.Open(s.storage.GetPath(ctx))
	s.NoError(err)
	defer file.Close()

	os.Remove(s.storage.GetPath(ctx))
	var storageData or.OrderStorage
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&storageData)
	s.NoError(err)

	storedOrder, exists := storageData.GetItem(ctx, order.ID)
	s.True(exists)
	s.Equal(order.UserID, storedOrder.UserID)
	s.Equal(order.State, storedOrder.State)
	s.Equal(order.KeepUntilDate.Unix(), storedOrder.KeepUntilDate.Unix())
}

func (s *OrderSuite) TestReadFromJSON() {
	ctx := context.Background()

	order := &models.Order{
		ID:            1,
		UserID:        1,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}

	s.storage.AddToStorage(ctx, order)

	err := s.storage.WriteToJSON()
	s.NoError(err)

	s.storage.Data = make(map[uint]*models.Order)

	err = s.storage.ReadFromJSON()
	s.NoError(err)

	os.Remove(s.storage.GetPath(ctx))

	storedOrder, exists := s.storage.GetItem(ctx, order.ID)
	s.True(exists)
	s.Equal(order.UserID, storedOrder.UserID)
	s.Equal(order.State, storedOrder.State)
	s.Equal(order.KeepUntilDate.Unix(), storedOrder.KeepUntilDate.Unix())
}
