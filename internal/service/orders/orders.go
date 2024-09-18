package orders

import (
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework-1.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/storage"
	"sort"
)

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
