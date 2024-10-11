package orders

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	"sort"
)

func AcceptOrder(ctx context.Context, s storage.Facade, or *models.Order) error {
	return s.AcceptOrder(ctx, or)
}

func PlaceOrder(ctx context.Context, s storage.Facade, ids []uint32) error {
	if len(ids) == 0 {
		return fmt.Errorf("Length of ids array is 0 ")
	}

	return s.PlaceOrder(ctx, ids)
}

func ReturnOrder(ctx context.Context, s storage.Facade, id uint) error {
	return s.ReturnOrder(ctx, id)
}

func ListOrders(ctx context.Context, s storage.Facade, id uint, n int, inPuP bool) error {
	list, err := s.ListOrders(ctx, id, inPuP)
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
