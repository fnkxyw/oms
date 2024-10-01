package orderStorage_test

import (
	"context"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/inmemory/orderStorage"
	"reflect"
	"sort"
	"testing"
	"time"
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
	ctx := context.Background()

	order := &models.Order{
		ID:     1,
		UserID: 100,
	}

	os.AddToStorage(ctx, order)

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

	ctx := context.Background()

	os := orderStorage.NewOrderStorage()
	os.AddToStorage(ctx, order)

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
			if got := os.IsConsist(ctx, tt.id); got != tt.want {
				t.Errorf("IsConsist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderStorage_DeleteOrderFromStorage(t *testing.T) {
	order := &models.Order{
		ID: 1,
	}

	ctx := context.Background()

	os := orderStorage.NewOrderStorage()
	os.AddToStorage(ctx, order)

	os.DeleteFromStorage(ctx, 1)

	if _, ok := os.Data[1]; ok {
		t.Errorf("DeleteOrderFromStorage() did not delete the order")
	}
}

func TestOrderStorage_GetOrder(t *testing.T) {
	order := &models.Order{ID: 1}

	ctx := context.Background()

	os := orderStorage.NewOrderStorage()
	os.AddToStorage(ctx, order)

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
			got, found := os.GetItem(ctx, tt.id)
			if !reflect.DeepEqual(got, tt.want) || found != tt.found {
				t.Errorf("GetOrder() got = %v, want %v, found = %v, wantFound = %v", got, tt.want, found, tt.found)
			}
		})
	}
}

func TestOrderStorage_GetOrderIDs(t *testing.T) {
	os := orderStorage.NewOrderStorage()
	ctx := context.Background()

	order1 := &models.Order{ID: 1}
	order2 := &models.Order{ID: 2}
	order3 := &models.Order{ID: 3}
	order4 := &models.Order{ID: 4}
	order5 := &models.Order{ID: 5}
	order6 := &models.Order{ID: 6}
	order7 := &models.Order{ID: 7}
	os.AddToStorage(ctx, order1)
	os.AddToStorage(ctx, order2)
	os.AddToStorage(ctx, order3)
	os.AddToStorage(ctx, order4)
	os.AddToStorage(ctx, order5)
	os.AddToStorage(ctx, order6)
	os.AddToStorage(ctx, order7)

	wantIDs := []uint{1, 2, 3, 4, 5, 6, 7}
	got, _ := os.GetIDs(ctx)
	time.Sleep(1 * time.Millisecond)
	sort.Slice(got, func(i, j int) bool {
		return got[i] < got[j]
	})

	if !reflect.DeepEqual(got, wantIDs) {
		t.Errorf("GetOrderIDs() = %v, want %v", got, wantIDs)
	}
}
