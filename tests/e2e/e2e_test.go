package e2e

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	c "gitlab.ozon.dev/akugnerevich/homework.git/internal/cli"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	"log"
	"os"
	"testing"
	"time"
)

const psqlDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

func TestApp(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	storageFacade := storage.NewStorageFacade(pool)

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

	if err = c.Run(ctx, storageFacade); err != nil {
		t.Errorf("Ошибка в e2e тесте: %v", err)
	}
	order, _ := storageFacade.GetItem(ctx, 10)
	if order.State != models.SoftDelete {
		t.Errorf("Ошибка: заказ с ID %d найден в orderStorage", 10)
	}
}

func TestAppSecond(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	storageFacade := storage.NewStorageFacade(pool)

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

	if err = c.Run(ctx, storageFacade); err != nil {
		t.Errorf("Ошибка в e2e тесте: %v", err)
	}

	order, ok := storageFacade.GetItem(ctx, 10)
	if !ok {
		t.Errorf("Ошибка: заказ с ID %d не найден в orderStorage", 10)
	}
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
		time.Sleep(120 * time.Millisecond)
	}
}
