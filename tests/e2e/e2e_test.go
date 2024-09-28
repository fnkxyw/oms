package e2e

import (
	"context"
	"fmt"
	c "gitlab.ozon.dev/akugnerevich/homework.git/internal/cli"
	orderStorage "gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/inmemory/orderStorage"
	"os"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	orderStorage := orderStorage.NewOrderStorage()
	ctx := context.Background()

	orderStorage.SetPath(ctx, "e2e_order.json")
	defer os.Remove(orderStorage.GetPath(ctx))

	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Ошибка при создании Pipe: %v", err)
	}
	os.Stdin = reader

	go func() {
		defer writer.Close()
		executeCommands(writer, []string{
			"acceptOrder\n",
			"10 10 2024-12-12\n",
			"10 999 box\n",
			"y\n",
			"listOrders\n", "10\n", "1\n",
			"placeOrder\n", "10\n",
			"listOrders\n", "10\n", "2\n", "1\n",
			"refundOrder\n", "10 10\n",
			"returnOrder\n", "10\n",
			"listReturns\n", "1 1\n",
			"exit\n",
		})
	}()

	if err = c.Run(ctx, orderStorage); err != nil {
		t.Errorf("Ошибка в e2e тесте: %v", err)
	}
	ok := orderStorage.IsConsist(ctx, 10)
	if ok {
		t.Errorf("Ошибка: заказ с ID %d найден в orderStorage", 10)
	}
}

func TestAppSecond(t *testing.T) {
	orderStorage := orderStorage.NewOrderStorage()
	ctx := context.Background()

	orderStorage.SetPath(ctx, "e2e_order2.json")
	defer os.Remove(orderStorage.GetPath(ctx))

	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Ошибка при создании Pipe: %v", err)
	}
	os.Stdin = reader

	go func() {
		defer writer.Close()
		executeCommands(writer, []string{
			"acceptOrder\n",
			"20 15 2024-12-12\n",
			"20 500 wrap\n",
			"n\n",
			"acceptOrder\n",
			"10 10 2024-12-12\n",
			"10 15 box\n",
			"y\n",
			"listOrders\n", "15\n", "1\n",
			"placeOrder\n", "20\n",
			"listOrders\n", "15\n", "2\n", "1\n",
			"refundOrder\n", "20 15\n",
			"exit\n",
		})
	}()

	if err = c.Run(ctx, orderStorage); err != nil {
		t.Errorf("Ошибка в e2e тесте: %v", err)
	}

	if !orderStorage.IsConsist(ctx, 10) {
		t.Errorf("Ошибка: заказ с ID %d не найден в orderStorage", 10)
	}
	order, _ := orderStorage.GetItem(ctx, 10)
	if order.Weight != 10 {
		t.Errorf("Данные в orderStorage неверные")
	}
}

func executeCommands(writer *os.File, commands []string) {
	for _, cmd := range commands {
		_, err := writer.Write([]byte(cmd))
		if err != nil {
			fmt.Printf("Ошибка записи команды %s: %v\n", cmd, err)
		}
		time.Sleep(1 * time.Millisecond)
	}
}
