package service

import (
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/storage"
	"time"
)

func WAcceptOrder(s *storage.OrderStorage) error {
	var order models.Order
	fmt.Println("Input OrderID _ CustomerID _ Date(form[2024-12(d)-12])")

	var dateString string
	_, err := fmt.Scan(&order.ID, &order.UserID, &dateString)
	if err != nil {
		return fmt.Errorf("Input data error", err)
	}

	order.Date, err = time.Parse("2006-01-02", dateString)
	if err != nil {
		return fmt.Errorf("Date parse error", err)
	}
	err = AcceptOrder(s, &order)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
