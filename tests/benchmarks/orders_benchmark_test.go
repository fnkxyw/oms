package benchmarks

import (
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	o "gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/inmemory/orderStorage"
	"os"
	"testing"
)

func BenchmarkAddOrderToStorage(b *testing.B) {
	orderStorage := o.NewOrderStorage()
	order := &models.Order{ID: 1}

	for i := 0; i < b.N; i++ {
		orderStorage.AddToStorage(order)
	}
}

func BenchmarkIsConsist(b *testing.B) {
	orderStorage := o.NewOrderStorage()
	orderStorage.AddToStorage(&models.Order{ID: 1})
	for i := 0; i < b.N; i++ {
		orderStorage.IsConsist(1)
	}
}

func BenchmarkDeleteOrderFromStorage(b *testing.B) {
	orderStorage := o.NewOrderStorage()
	order := &models.Order{ID: 1}
	orderStorage.AddToStorage(order)

	for i := 0; i < b.N; i++ {
		orderStorage.DeleteFromStorage(1)
		orderStorage.AddToStorage(order)
	}
}

func BenchmarkGetOrder(b *testing.B) {
	orderStorage := o.NewOrderStorage()
	order := &models.Order{ID: 1}
	orderStorage.AddToStorage(order)

	for i := 0; i < b.N; i++ {
		orderStorage.GetItem(1)
	}
}

func BenchmarkGetOrderIDs(b *testing.B) {
	orderStorage := o.NewOrderStorage()
	for i := uint(0); i < 100; i++ {
		orderStorage.AddToStorage(&models.Order{ID: i})
	}

	for i := 0; i < b.N; i++ {
		orderStorage.GetIDs()
	}
}

func BenchmarkReadFromJSON(b *testing.B) {
	orderStorage := o.NewOrderStorage()
	orderStorage.SetPath("orders_test_benchmark.json")
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

	os.Remove(orderStorage.GetPath())
}

func BenchmarkWriteToJSON(b *testing.B) {
	orderStorage := o.NewOrderStorage()
	orderStorage.SetPath("orders_test_benchmark.json")
	for i := uint(0); i < 100; i++ {
		orderStorage.AddToStorage(&models.Order{ID: i})
	}

	for i := 0; i < b.N; i++ {
		err := orderStorage.WriteToJSON()
		if err != nil {
			b.Fatal(err)
		}
	}
	os.Remove(orderStorage.GetPath())

}
