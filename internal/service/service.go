package service

import (
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/storage"
	"time"
)

func AcceptOrder(s *storage.OrderStorage, or *models.Order) error {
	if or.Date.Before(time.Now()) {
		return fmt.Errorf("Incorrect date")
	}
	or.Accept = true
	err := s.AddOrderToStorage(or)
	if err != nil {
		return err
	}
	fmt.Println("Correct!")
	return nil
}
