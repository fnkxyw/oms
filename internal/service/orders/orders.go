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

func FilterOrders(s *storage.OrderStorage, id uint, inPuP bool) []*models.Order {
	var filtered []*models.Order

	for _, order := range s.Data {
		if order.UserID == id && (!inPuP || (order.State == models.AcceptState || order.State == models.ReturnedState)) {
			filtered = append(filtered, order)
		}
	}

	return filtered
}

func CheckIDsOrders(s *storage.OrderStorage, ids []uint) error {
	temp := s.Data[ids[0]].UserID
	for _, id := range ids {
		if s.Data[id].UserID != temp {
			return e.ErrorNotAllIDs
		}
	}
	return nil
}
