package controller

import (
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/controller/inputs"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders/packing"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/returns"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
)

func WAcceptOrder(s storage.Storage) error {
	order, packageType, needWrapping, err := inputs.CollectOrderInput()
	if err != nil {
		return err
	}

	err = packing.Packing(order, packageType, needWrapping)
	if err != nil {
		return err
	}

	err = orders.AcceptOrder(s, order)
	if err != nil {
		return err
	}

	return nil
}

func WReturnOrder(s storage.Storage) error {
	id, err := inputs.InputOrderID()
	if err != nil {
		return err
	}

	return orders.ReturnOrder(s, id)
}

func WPlaceOrder(s storage.Storage) error {
	uintdata, err := inputs.InputOrderIDs()
	if err != nil {
		return err
	}

	return orders.PlaceOrder(s, uintdata)
}

func WListOrders(s storage.Storage) error {
	id, err := inputs.InputUserID()
	if err != nil {
		return err
	}

	temp, err := inputs.InputListChoice()
	if err != nil {
		return err
	}

	switch temp {
	case 1:
		return orders.ListOrders(s, id, 0, true)
	case 2:
		n, err := inputs.InputN()
		if err != nil {
			return err
		}
		return orders.ListOrders(s, id, n, false)
	}

	return nil
}

func WRefundOrder(oS storage.Storage) error {
	orderId, userId, err := inputs.InputOrderAndUserID()
	if err != nil {
		return err
	}

	return returns.RefundOrder(oS, orderId, userId)
}

func WListReturns(oS storage.Storage) error {
	limit, page, err := inputs.InputReturnsPagination()
	if err != nil {
		return err
	}

	return returns.ListReturns(oS, limit, page)
}
