package service

import (
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework-1.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/service/orders"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/storage"
	"time"
)

// принять заказ от курьера
func AcceptOrder(s *storage.OrderStorage, or *models.Order) error {
	if s.IsConsist(or.ID) {
		return e.ErrIsConsist
	}
	if or.KeepUntilDate.Before(time.Now()) {
		return e.ErrDate
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
	return s.Data[id].CanReturned()
}

// доставить заказ юзеру
func PlaceOrder(s *storage.OrderStorage, ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("Length of ids array is 0 ")
	}

	orders.CheckIDsOrders(s, ids)

	for _, id := range ids {
		order := s.Data[id]

		if order.State == models.PlaceState {
			return fmt.Errorf("Order by id: %d is already placed", id)
		}

		if order.State == models.SoftDelete {
			return fmt.Errorf("Order by id: %d was deleted", id)
		}

		if !order.KeepUntilDate.After(time.Now()) {
			return fmt.Errorf("Order by id: %d cannot be issued to the customer because the date is invalid", id)
		}

		order.State = models.PlaceState
		order.PlaceDate = time.Now()
	}

	return nil
}

func ListOrders(s *storage.OrderStorage, id uint, n int, inPuP bool) error {
	var list []*models.Order
	list = orders.FilterOrders(s, id, inPuP)
	orders.SortOrders(list)
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
func RefundOrder(rs *storage.ReturnStorage, os *storage.OrderStorage, id uint, userId uint) error {
	if !os.IsConsist(id) {
		return e.ErrCheckOrderID
	}
	if os.Data[id].State != models.PlaceState {
		return e.ErrNotPlace
	}
	if time.Now().After(os.Data[id].PlaceDate.AddDate(0, 0, 2)) {
		return e.ErrTimeExpired
	}
	if os.Data[id].UserID != userId {
		return e.ErrIncorrectUserId
	}

	rs.AddReturnToStorage(&models.Return{
		ID:           id,
		UserID:       userId,
		DateOfReturn: time.Now(),
	})
	os.Data[id].State = models.ReturnedState

	return nil
}

// показать список возвратов с    пагинацией
func ListReturns(rs *storage.ReturnStorage, limit, page int) error {
	var list []*models.Return
	for _, v := range rs.Data {
		list = append(list, v)
	}
	return pagePagination(list, page, limit)
}
