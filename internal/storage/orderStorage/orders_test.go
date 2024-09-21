package orderStorage_test

import (
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/orderStorage"
	"reflect"
	"sort"
	"testing"
)

func TestNewOrderStorage(t *testing.T) {
	want := &orderStorage.OrderStorage{Data: make(map[uint]*models.Order)}
	got := orderStorage.NewOrderStorage()
	if !reflect.DeepEqual(got.Data, want.Data) {
		t.Errorf("NewOrderStorage() = %v, want %v", got, want)
	}
}

func TestOrderStorage_AddOrderToStorage(t *testing.T) {
	os := orderStorage.NewOrderStorage()

	order := &models.Order{
		ID:     1,
		UserID: 100,
	}

	os.AddOrderToStorage(order)

	if len(os.Data) != 1 {
		t.Errorf("expected 1 order, got %d", len(os.Data))
	}

	if os.Data[1] != order {
		t.Errorf("order was not correctly added to storage")
	}
}

func TestOrderStorage_IsConsist(t *testing.T) {
	order := &models.Order{
		ID: 1,
	}

	os := orderStorage.NewOrderStorage()
	os.AddOrderToStorage(order)

	tests := []struct {
		name string
		id   uint
		want bool
	}{
		{
			name: "Order exists",
			id:   1,
			want: true,
		},
		{
			name: "Order does not exist",
			id:   2,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := os.IsConsist(tt.id); got != tt.want {
				t.Errorf("IsConsist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderStorage_DeleteOrderFromStorage(t *testing.T) {
	order := &models.Order{
		ID: 1,
	}

	os := orderStorage.NewOrderStorage()
	os.AddOrderToStorage(order)

	os.DeleteOrderFromStorage(1)

	if _, ok := os.Data[1]; ok {
		t.Errorf("DeleteOrderFromStorage() did not delete the order")
	}
}

func TestOrderStorage_GetOrder(t *testing.T) {
	order := &models.Order{ID: 1}

	os := orderStorage.NewOrderStorage()
	os.AddOrderToStorage(order)

	tests := []struct {
		name  string
		id    uint
		want  *models.Order
		found bool
	}{
		{
			name:  "Order found",
			id:    1,
			want:  order,
			found: true,
		},
		{
			name:  "Order not found",
			id:    2,
			want:  nil,
			found: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := os.GetOrder(tt.id)
			if !reflect.DeepEqual(got, tt.want) || found != tt.found {
				t.Errorf("GetOrder() got = %v, want %v, found = %v, wantFound = %v", got, tt.want, found, tt.found)
			}
		})
	}
}

func TestOrderStorage_GetOrderIDs(t *testing.T) {
	os := orderStorage.NewOrderStorage()

	order1 := &models.Order{ID: 1}
	order2 := &models.Order{ID: 2}
	order3 := &models.Order{ID: 3}
	os.AddOrderToStorage(order1)
	os.AddOrderToStorage(order2)
	os.AddOrderToStorage(order3)

	wantIDs := []uint{1, 2, 3}
	gotIDs := os.GetOrderIDs()

	sort.SliceIsSorted(gotIDs, func(i, j int) bool {
		return i < j
	})

	if !reflect.DeepEqual(gotIDs, wantIDs) {
		t.Errorf("GetOrderIDs() = %v, want %v", gotIDs, wantIDs)
	}
}
