package orders

import (
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework-1.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/service/pagination"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/storage"
	"sort"
	"time"
)

// принять заказ от курьера
func AcceptOrder(s storage.OrderStorageInterface, or *models.Order) error {
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

// доставить заказ юзеру
func PlaceOrder(s storage.OrderStorageInterface, ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("Length of ids array is 0 ")
	}

	CheckIDsOrders(s, ids)

	for _, id := range ids {
		order, exists := s.GetOrder(id)
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

func ListOrders(s storage.OrderStorageInterface, id uint, n int, inPuP bool) error {
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

// вернуть заказ юзеру
func RefundOrder(rs storage.ReturnStorageInterface, os storage.OrderStorageInterface, id uint, userId uint) error {
	order, exists := os.GetOrder(id)
	if !exists {
		return e.ErrCheckOrderID
	}
	if order.State != models.PlaceState {
		return e.ErrNotPlace
	}
	if time.Now().After(order.PlaceDate.AddDate(0, 0, 2)) {
		return e.ErrTimeExpired
	}
	if order.UserID != userId {
		return e.ErrIncorrectUserId
	}

	rs.AddReturnToStorage(&models.Return{
		ID:           id,
		UserID:       userId,
		DateOfReturn: time.Now(),
	})
	order.State = models.ReturnedState

	return nil
}

func SortOrders(o []*models.Order) error {
	sort.Slice(o, func(i, j int) bool {
		return o[i].AcceptTime.Before(o[j].AcceptTime)
	})
	return nil
}

func FilterOrders(s storage.OrderStorageInterface, id uint, inPuP bool) []*models.Order {
	var filtered []*models.Order
	for _, o := range s.GetOrderIDs() {
		order, exists := s.GetOrder(o)
		if !exists {
			continue
		}
		if order.UserID == id && (!inPuP || (order.State == models.AcceptState || order.State == models.ReturnedState)) {
			filtered = append(filtered, order)
		}
	}

	return filtered
}

func CheckIDsOrders(s storage.OrderStorageInterface, ids []uint) error {
	order, _ := s.GetOrder(ids[0])
	temp := order.UserID
	for _, id := range ids {
		order, _ = s.GetOrder(id)
		if order.UserID != temp {
			return e.ErrNotAllIDs
		}
	}
	return nil
}
