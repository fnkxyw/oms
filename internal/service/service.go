package service

import (
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/storage"
	"sort"
	"time"
)

func AcceptOrder(s *storage.OrderStorage, or *models.Order) error {
	if or.Date.Before(time.Now()) {
		return fmt.Errorf("Incorrect date: %w")
	}
	or.Accept = true
	or.AcceptTime = time.Now()
	err := s.AddOrderToStorage(or)
	if err != nil {
		return err
	}
	fmt.Println("Correct!")
	return nil
}

func ReturnOrder(s *storage.OrderStorage, id uint) error {
	if s.Data[id].Issued == false && s.Data[id].Date.After(time.Now()) {
		s.DeleteOrderFromStorage(id)
	} else {
		return fmt.Errorf("Order can`t be returned")
	}
	return nil
}

func PlaceOrder(s *storage.OrderStorage, id []uint) error {
	if len(id) == 0 {
		return fmt.Errorf("Length of ids array is 0")
	}
	temp := s.Data[id[0]].UserID
	for _, v := range id {
		if s.Data[v].Issued == true {
			return fmt.Errorf("Order by id: %d already place\n", v)

		}
		if s.Data[v].UserID == temp && s.Data[v].Date.After(time.Now()) {
			s.Data[v].Issued = true
			s.Data[v].IssuedDate = time.Now()
		} else {
			return fmt.Errorf("You cannot issue this item to the customer")
		}
	}
	return nil
}

func ListOrders(s *storage.OrderStorage, id uint, n int, inPuP bool) error {
	var list []*models.Order
	for _, v := range s.Data {
		if inPuP == true {
			if v.UserID == id && v.Issued == false {
				list = append(list, v)
			}
		} else {
			if v.UserID == id {
				list = append(list, v)
			}
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return !(list[i].AcceptTime.Before(list[j].AcceptTime))
	})
	if n < 1 {
		n = 1
	}
	if inPuP == false {
		list = list[:n]
	}
	for _, order := range list {
		fmt.Printf("OrderID: %v, Reciver: %v, IssuedStatus: %v, Date until which it will be stored: %v \n", order.ID, order.UserID, order.Issued, order.Date)
	}
	return nil
}
