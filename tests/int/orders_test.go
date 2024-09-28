package int

import (
	"context"
	"encoding/json"
	"errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/inmemory/orderStorage"
	"os"
	"testing"
	"time"
)

func newStorage() *orderStorage.OrderStorage {
	storage := orderStorage.NewOrderStorage()
	ctx := context.Background()

	storage.SetPath(ctx, "order_test.json")
	return storage
}

func TestAcceptOrder(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	storage := newStorage()
	order := &models.Order{
		ID:            1,
		UserID:        1,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}

	err := orders.AcceptOrder(ctx, storage, order)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	storedOrder, exists := storage.GetItem(ctx, order.ID)
	if !exists {
		t.Fatal("expected order to exist in storage")
	}
	if storedOrder.State != models.AcceptState {
		t.Errorf("expected order state %v, got %v", models.AcceptState, storedOrder.State)
	}
}

func TestAcceptOrder_PastDate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	storage := newStorage()
	order := &models.Order{
		ID:            3,
		UserID:        1,
		KeepUntilDate: time.Now().Add(-1 * time.Hour),
	}

	err := orders.AcceptOrder(ctx, storage, order)
	if !errors.Is(err, e.ErrDate) {
		t.Errorf("expected error %v, got %v", e.ErrDate, err)
	}
}

func TestAcceptOrder_EqualOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	storage := newStorage()
	order1 := &models.Order{
		ID:            2,
		UserID:        1,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}
	order2 := &models.Order{
		ID:            2,
		UserID:        1,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}

	err := orders.AcceptOrder(ctx, storage, order1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = orders.AcceptOrder(ctx, storage, order2)
	if !errors.Is(err, e.ErrIsConsist) {
		t.Errorf("expected error %v, got %v", e.ErrIsConsist, err)
	}
}

func TestPlaceOrder(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	storage := newStorage()
	order := &models.Order{
		ID:            2,
		UserID:        1,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
		State:         models.AcceptState,
	}
	storage.AddToStorage(ctx, order)

	err := orders.PlaceOrder(ctx, storage, []uint{2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedOrder, exists := storage.GetItem(ctx, 2)
	if !exists {
		t.Fatal("expected order to exist in storage")
	}
	if updatedOrder.State != models.PlaceState {
		t.Errorf("expected order state %v, got %v", models.PlaceState, updatedOrder.State)
	}
}

func TestPlaceOrder_NoConsist(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	storage := newStorage()
	err := orders.PlaceOrder(ctx, storage, []uint{1})
	if !errors.Is(err, e.ErrNoConsist) {
		t.Errorf("expected error %v, got %v", e.ErrNoConsist, err)
	}
}

func TestListOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	storage := newStorage()
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
	storage.AddToStorage(ctx, order1)
	storage.AddToStorage(ctx, order2)

	err := orders.ListOrders(ctx, storage, 1, 2, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	id := storage.GetIDs(ctx)
	if len(id) != 2 {
		t.Errorf("expected 2 orders, got %d", len(id))
	}
	if !contain(id, 4) || !contain(id, 5) {
		t.Error("expected orders with ID 4 and 5")
	}
}

func TestReturnOrder(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	storage := newStorage()
	order := &models.Order{
		ID:     6,
		UserID: 1,
		State:  models.RefundedState,
	}
	storage.AddToStorage(ctx, order)

	err := orders.ReturnOrder(ctx, storage, 6)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.CanReturned() != nil {
		t.Error("expected order to be not returnable")
	}
}

func TestReturnOrder_NoConsist(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	storage := newStorage()
	order := &models.Order{
		ID:     1,
		UserID: 1,
		State:  models.RefundedState,
	}
	err := orders.ReturnOrder(ctx, storage, order.ID)
	if !errors.Is(err, e.ErrNoConsist) {
		t.Errorf("expected error %v, got %v", e.ErrNoConsist, err)
	}
}

func TestCheckIDOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	storage := newStorage()
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
	storage.AddToStorage(ctx, order1)
	storage.AddToStorage(ctx, order2)
	storage.AddToStorage(ctx, order3)

	err := orders.CheckIDsOrders(ctx, storage, []uint{10, 11})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = orders.CheckIDsOrders(ctx, storage, []uint{10, 12})
	if !errors.Is(err, e.ErrNotAllIDs) {
		t.Errorf("expected error %v, got %v", e.ErrNotAllIDs, err)
	}
}

func TestWriteToJSON(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	storage := newStorage()
	order := &models.Order{
		ID:            1,
		UserID:        1,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}

	storage.AddToStorage(ctx, order)

	err := storage.WriteToJSON()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	file, err := os.Open(storage.GetPath(ctx))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer file.Close()
	defer os.Remove(storage.GetPath(ctx))

	var storageData orderStorage.OrderStorage
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&storageData)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	storedOrder, exists := storageData.GetItem(ctx, order.ID)
	if !exists {
		t.Fatal("expected order to exist in storage")
	}
	if order.UserID != storedOrder.UserID || order.State != storedOrder.State || order.KeepUntilDate.Unix() != storedOrder.KeepUntilDate.Unix() {
		t.Error("order data does not match")
	}
}

func TestReadFromJSON(t *testing.T) {
	ctx := context.Background()

	storage := newStorage()
	order := &models.Order{
		ID:            1,
		UserID:        1,
		State:         models.AcceptState,
		KeepUntilDate: time.Now().Add(24 * time.Hour),
	}

	storage.AddToStorage(ctx, order)

	err := storage.WriteToJSON()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	storage.Data = make(map[uint]*models.Order)

	err = storage.ReadFromJSON()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer os.Remove(storage.GetPath(ctx))

	storedOrder, exists := storage.GetItem(ctx, order.ID)
	if !exists {
		t.Fatal("expected order to exist in storage")
	}
	if order.UserID != storedOrder.UserID || order.State != storedOrder.State || order.KeepUntilDate.Unix() != storedOrder.KeepUntilDate.Unix() {
		t.Error("order data does not match")
	}
}

func contain(ids []uint, id uint) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}
