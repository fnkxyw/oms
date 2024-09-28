package benchmarks

import (
	"context"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	o "gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/inmemory/orderStorage"
	"os"
	"testing"
)

func BenchmarkAddOrderToStorage(b *testing.B) {
	orderStorage := o.NewOrderStorage()
	ctx := context.Background()
	order := &models.Order{ID: 1}

	for i := 0; i < b.N; i++ {
		orderStorage.AddToStorage(ctx, order)
	}
}

func BenchmarkIsConsist(b *testing.B) {
	orderStorage := o.NewOrderStorage()
	ctx := context.Background()
	orderStorage.AddToStorage(ctx, &models.Order{ID: 1})
	for i := 0; i < b.N; i++ {
		orderStorage.IsConsist(ctx, 1)
	}
}

func BenchmarkDeleteOrderFromStorage(b *testing.B) {
	orderStorage := o.NewOrderStorage()
	order := &models.Order{ID: 1}
	ctx := context.Background()

	orderStorage.AddToStorage(ctx, order)

	for i := 0; i < b.N; i++ {
		orderStorage.DeleteFromStorage(ctx, 1)
		orderStorage.AddToStorage(ctx, order)
	}
}

func BenchmarkGetOrder(b *testing.B) {
	orderStorage := o.NewOrderStorage()
	order := &models.Order{ID: 1}
	ctx := context.Background()

	orderStorage.AddToStorage(ctx, order)

	for i := 0; i < b.N; i++ {
		orderStorage.GetItem(ctx, 1)
	}
}

func BenchmarkGetOrderIDs(b *testing.B) {
	orderStorage := o.NewOrderStorage()
	ctx := context.Background()

	for i := uint(0); i < 100; i++ {
		orderStorage.AddToStorage(ctx, &models.Order{ID: i})
	}

	for i := 0; i < b.N; i++ {
		orderStorage.GetIDs(ctx)
	}
}

func BenchmarkReadFromJSON(b *testing.B) {
	orderStorage := o.NewOrderStorage()
	ctx := context.Background()

	orderStorage.SetPath(ctx, "orders_test_benchmark.json")
	err := orderStorage.WriteToJSON()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		err := orderStorage.ReadFromJSON()
		if err != nil {
			b.Fatal(err)
		}
	}

	os.Remove(orderStorage.GetPath(ctx))
}

func BenchmarkWriteToJSON(b *testing.B) {
	orderStorage := o.NewOrderStorage()
	ctx := context.Background()

	orderStorage.SetPath(ctx, "orders_test_benchmark.json")
	for i := uint(0); i < 100; i++ {
		orderStorage.AddToStorage(ctx, &models.Order{ID: i})
	}

	for i := 0; i < b.N; i++ {
		err := orderStorage.WriteToJSON()
		if err != nil {
			b.Fatal(err)
		}
	}
	os.Remove(orderStorage.GetPath(ctx))

}
