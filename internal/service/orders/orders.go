package orders

import (
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/pagination"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	"sort"
	"time"
)

func AcceptOrder(s storage.Storage, or *models.Order) error {
	if s.IsConsist(or.ID) {
		return e.ErrIsConsist
	}
	if or.KeepUntilDate.Before(time.Now()) {
		return e.ErrDate
	}
	or.State = models.AcceptState
	or.AcceptTime = time.Now().Unix()
	s.AddToStorage(or)
	return nil
}

// доставить заказ юзеру
func PlaceOrder(s storage.Storage, ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("Length of ids array is 0 ")
	}

	err := CheckIDsOrders(s, ids)
	if err != nil {
		return err
	}

	for _, id := range ids {
		order, exists := s.GetItem(id)
		if !exists {
			return e.ErrNoConsist
		}
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

func ListOrders(s storage.Storage, id uint, n int, inPuP bool) error {
	var list []*models.Order
	list = FilterOrders(s, id, inPuP)
	SortOrders(list)
	if n < 1 {
		n = 1
	} else if n > len(list) {
		n = len(list)
	}
	if !inPuP {
		list = list[:n]
	}
	return pagination.ScrollPagination(list, 1)
}

// вернуть заказ курьеру
func ReturnOrder(s storage.Storage, id uint) error {

	order, exists := s.GetItem(id)
	if !exists {
		return e.ErrNoConsist
	}
	return order.CanReturned()
}

func SortOrders(o []*models.Order) error {
	sort.Slice(o, func(i, j int) bool {
		return o[i].AcceptTime < o[j].AcceptTime
	})
	return nil
}

func FilterOrders(s storage.Storage, id uint, inPuP bool) []*models.Order {
	var filtered []*models.Order
	var ids []uint
	ids = s.GetIDs()
	for _, o := range ids {
		order, exists := s.GetItem(o)
		if exists && order.UserID == id && (!inPuP || order.State == models.AcceptState || order.State == models.ReturnedState) {
			filtered = append(filtered, order)
		}
	}

	return filtered
}

func CheckIDsOrders(s storage.Storage, ids []uint) error {
	order, ok := s.GetItem(ids[0])
	if !ok {
		return e.ErrNoConsist
	}
	temp := order.UserID
	for _, id := range ids {
		order, _ = s.GetItem(id)
		if order.UserID != temp {
			return e.ErrNotAllIDs
		}
	}
	return nil
}
