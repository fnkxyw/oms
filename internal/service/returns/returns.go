package returns

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"

	"time"
)

func RefundOrder(ctx context.Context, os storage.Storage, id uint, userId uint) error {
	order, exists := os.GetItem(ctx, id)

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

	err := os.UpdateState(ctx, id, models.RefundedState)
	if err != nil {
		return err
	}

	return nil
}

func ListReturns(ctx context.Context, os storage.Storage, limit, page int) error {
	var list []models.Order
	list, err := os.GetReturns(ctx, page-1, limit)
	if err != nil {
		return err
	}
	for _, v := range list {
		fmt.Printf("OrderID: %d, UserID: %d \n", v.ID, v.UserID)
	}
	return nil
}
