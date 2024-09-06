package service

import (
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/storage"
	"sort"
	"time"
)

// принять заказ от курьера
func AcceptOrder(s *storage.OrderStorage, or *models.Order) error {
	if or.Date.Before(time.Now()) {
		return fmt.Errorf("Incorrect date \n")
	}
	or.State = models.AcceptState
	or.AcceptTime = time.Now()
	err := s.AddOrderToStorage(or)
	if err != nil {
		return err
	}
	fmt.Println("Correct!")
	return nil
}

// вернуть заказ курьеру
func ReturnOrder(s *storage.OrderStorage, id uint) error {
	if s.Data[id].State == models.ReturnedState || s.Data[id].Date.Before(time.Now()) {
		s.DeleteOrderFromStorage(id)
	} else {
		return fmt.Errorf("Order can`t be returned\n")
	}
	return nil
}

// доставить заказ юзеру
func PlaceOrder(s *storage.OrderStorage, id []uint) error {
	if len(id) == 0 {
		return fmt.Errorf("Length of ids array is 0\n")
	}
	temp := s.Data[id[0]].UserID
	for _, v := range id {
		if s.Data[v].UserID != temp {
			return fmt.Errorf("Not all orders for one user \n")
		}
	}
	for _, v := range id {
		if s.Data[v].State == models.PlaceState {
			return fmt.Errorf("Order by id: %d already place\n", v)

		}
		if s.Data[v].Date.After(time.Now()) {
			s.Data[v].State = models.PlaceState
			s.Data[v].PlaceData = time.Now()
		} else {
			return fmt.Errorf("You cannot issue this item to the customer\n")
		}
	}
	return nil
}

// /показать список заказов
func ListOrders(s *storage.OrderStorage, id uint, n int, inPuP bool) error {
	var list []*models.Order
	for _, v := range s.Data {
		if inPuP == true {
			if v.UserID == id && (v.State == models.AcceptState || v.State == models.ReturnedState) {
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
	} else if n > len(list) {
		n = len(list)
	}
	if inPuP == false {
		list = list[:n]
	}
	return scrollPagination(list, 1)
}

// вернуть заказ юзеру
func ReturnUser(rs *storage.ReturnStorage, os *storage.OrderStorage, id uint, userId uint) error {
	if !os.IsConsist(id) {
		return fmt.Errorf("Check input OrderId\n")
	}
	if os.Data[id].State != models.PlaceState {
		return fmt.Errorf("Order are not placed \n")
	}
	if time.Now().After(os.Data[id].PlaceData.AddDate(0, 0, 2)) {
		return fmt.Errorf("Return time has expired :( \n")
	}
	if os.IsConsist(id) && os.Data[id].UserID == userId && os.Data[id].State == models.PlaceState {
		rs.AddReturnToStorage(&models.Return{
			ID:           id,
			UserID:       userId,
			DateOfReturn: time.Now(),
		})
		os.Data[id].State = models.ReturnedState
	} else {
		return fmt.Errorf("Check input data\n")
	}
	return nil
}

// показать список возвратов с постраничной пагинацией
func ListReturns(rs *storage.ReturnStorage, limit, page int) error {
	var list []*models.Return
	for _, v := range rs.Data {
		list = append(list, v)
	}
	return pagePagination(list, page, limit)
}
