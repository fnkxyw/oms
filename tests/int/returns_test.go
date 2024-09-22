package int

import (
	"encoding/json"
	"errors"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/returns"
	"os"
	"testing"
	"time"

	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/orderStorage"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/returnStorage"
)

func newOrderStorage() *orderStorage.OrderStorage {
	return orderStorage.NewOrderStorage()
}

func newReturnStorage() *returnStorage.ReturnStorage {
	r := returnStorage.NewReturnStorage()
	r.SetPath("returns_test.json")
	return r
}

func TestRefundOrder(t *testing.T) {
	t.Parallel()
	orderStorage := newOrderStorage()
	returnStorage := newReturnStorage()

	order := &models.Order{
		ID:        1,
		UserID:    1,
		State:     models.PlaceState,
		PlaceDate: time.Now(),
	}
	orderStorage.AddOrderToStorage(order)

	err := returns.RefundOrder(returnStorage, orderStorage, 1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.State != models.ReturnedState {
		t.Errorf("expected order state %v, got %v", models.ReturnedState, order.State)
	}

	if !returnStorage.IsConsist(1) {
		t.Error("expected return to exist in storage")
	}
}

func TestRefundOrder_NoOrder(t *testing.T) {
	t.Parallel()
	orderStorage := newOrderStorage()
	returnStorage := newReturnStorage()

	err := returns.RefundOrder(returnStorage, orderStorage, 1, 1)
	if !errors.Is(err, e.ErrCheckOrderID) {
		t.Errorf("expected error %v, got %v", e.ErrCheckOrderID, err)
	}
}

func TestListReturn(t *testing.T) {
	t.Parallel()
	returnStorage := newReturnStorage()

	returnStorage.AddReturnToStorage(&models.Return{ID: 1, UserID: 1})

	err := returns.ListReturns(returnStorage, 10, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWriteToJSON_ReturnStorage(t *testing.T) {
	storage := newReturnStorage()

	returnItem := &models.Return{
		ID:     1,
		UserID: 1,
	}

	storage.AddReturnToStorage(returnItem)
	storage.SetPath("returns_test.json")
	err := storage.WriteToJSON()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	file, err := os.Open(storage.GetPath())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer file.Close()
	defer os.Remove(storage.GetPath())

	var storageData returnStorage.ReturnStorage
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&storageData)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	storedReturn, exists := storageData.GetReturn(returnItem.ID)
	if !exists {
		t.Fatal("expected return to exist in storage")
	}
	if returnItem.UserID != storedReturn.UserID {
		t.Error("return data does not match")
	}
}

func TestReadFromJSON_ReturnStorage(t *testing.T) {

	storage := newReturnStorage()
	returnItem := &models.Return{
		ID:     1,
		UserID: 1,
	}

	storage.AddReturnToStorage(returnItem)
	storage.SetPath("returns_test.json")

	err := storage.WriteToJSON()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	storage.Data = make(map[uint]*models.Return)

	err = storage.ReadFromJSON()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	defer os.Remove(storage.GetPath())

	storedReturn, exists := storage.GetReturn(returnItem.ID)
	if !exists {
		t.Fatal("expected return to exist in storage")
	}
	if returnItem.UserID != storedReturn.UserID {
		t.Error("return data does not match")
	}
}
