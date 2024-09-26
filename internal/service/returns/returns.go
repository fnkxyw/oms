package returns

import (
	"errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/pagination"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/orderStorage"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/returnStorage"
	"time"
)

func RefundOrder(rs returnStorage.ReturnStorageInterface, os orderStorage.OrderStorageInterface, id uint, userId uint) error {
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

	err := rs.AddReturnToStorage(&models.Return{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		return errors.New("add return to storage error in refund order")
	}
	order.State = models.ReturnedState

	return nil
}

func ListReturns(rs returnStorage.ReturnStorageInterface, limit, page int) error {
	var list []*models.Return
	for _, v := range rs.GetReturnIDs() {
		r, _ := rs.GetReturn(v)
		list = append(list, r)
	}
	return pagination.PagePagination(list, page, limit)
}
