package returns

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
)

func RefundOrder(ctx context.Context, os storage.Facade, id uint, userId uint) error {
	return os.RefundOrder(ctx, id, userId)
}

func ListReturns(ctx context.Context, os storage.Facade, limit, page int) error {
	list, err := os.ListReturns(ctx, limit, page)
	if err != nil {
		return err
	}
	for _, v := range list {
		fmt.Printf("OrderID: %d, UserID: %d \n", v.ID, v.UserID)
	}
	return nil
}
