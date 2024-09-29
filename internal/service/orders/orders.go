package orders

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	"sort"
	"time"
)

func AcceptOrder(ctx context.Context, s storage.Storage, or *models.Order) error {
	if s.IsConsist(ctx, or.ID) {
		return e.ErrIsConsist
	}
	if or.KeepUntilDate.Before(time.Now()) {
		return e.ErrDate
	}
	or.State = models.AcceptState
	or.AcceptTime = time.Now().Unix()
	err := s.AddToStorage(ctx, or)
	if err != nil {
		return err
	}
	return nil
}

// доставить заказ юзеру
func PlaceOrder(ctx context.Context, s storage.Storage, ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("Length of ids array is 0 ")
	}

	err := CheckIDsOrders(ctx, s, ids)
	if err != nil {
		return err
	}

	for _, id := range ids {
		order, exists := s.GetItem(ctx, id)
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

		err := s.UpdateBeforePlace(ctx, id, models.PlaceState, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

// вернуть заказ курьеру
func ReturnOrder(ctx context.Context, s storage.Storage, id uint) error {

	order, exists := s.GetItem(ctx, id)
	if !exists {
		return e.ErrNoConsist
	}
	err := order.CanReturned()
	if err != nil {
		return err
	}

	err = s.DeleteFromStorage(ctx, order.ID)
	if err != nil {
		return err
	}

	return nil
}

func ListOrders(ctx context.Context, s storage.Storage, id uint, n int, inPuP bool) error {
	var list []models.Order
	list, err := s.GetOrders(ctx, id, inPuP)
	if err != nil {
		return err
	}
	SortOrders(list)
	if n < 1 {
		n = 1
	} else if n > len(list) {
		n = len(list)
	}
	if !inPuP {
		list = list[:n]
	}
	for _, v := range list {
		fmt.Printf("OrderID: %v, Reciver: %v, State: %v, Price: %v₽, Date until which it will be stored: %v \n",
			v.ID, v.UserID, v.State, v.Price, v.KeepUntilDate)
	}
	return nil
}

func SortOrders(o []models.Order) {
	sort.Slice(o, func(i, j int) bool {
		return o[i].AcceptTime > o[j].AcceptTime
	})
}

func CheckIDsOrders(ctx context.Context, s storage.Storage, ids []uint) error {
	order, ok := s.GetItem(ctx, ids[0])
	if !ok {
		return e.ErrNoConsist
	}
	temp := order.UserID
	for _, id := range ids {
		order, _ = s.GetItem(ctx, id)
		if order.UserID != temp {
			return e.ErrNotAllIDs
		}
	}
	return nil
}
