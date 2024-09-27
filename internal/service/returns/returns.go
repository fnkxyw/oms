package returns

import (
	"context"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/pagination"
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

	os.UpdateState(ctx, id, models.RefundedState)

	return nil
}

func ListReturns(ctx context.Context, os storage.Storage, limit, page int) error {
	var list []*models.Order
	for _, v := range os.GetIDs(ctx) {
		order, _ := os.GetItem(ctx, v)
		if order.State == models.SoftDelete {
			list = append(list, order)
		}
	}
	return pagination.PagePagination(list, page, limit)
}
