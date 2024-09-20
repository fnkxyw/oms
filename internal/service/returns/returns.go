package returns

import (
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/pagination"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
)

func ReturnOrder(s storage.OrderStorageInterface, id uint) error {
	order, exists := s.GetOrder(id)
	if !exists {
		return e.ErrNoConsist
	}
	return order.CanReturned()
}

func ListReturns(rs storage.ReturnStorageInterface, limit, page int) error {
	var list []*models.Return
	for _, v := range rs.GetReturnIDs() {
		r, _ := rs.GetReturn(v)
		list = append(list, r)
	}
	return pagination.PagePagination(list, page, limit)
}
