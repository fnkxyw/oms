package controller

import (
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/controller/inputs"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders/packing"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/returns"
	s "gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/orderStorage"
	r "gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/returnStorage"
)

func WAcceptOrder(s s.OrderStorageInterface) error {
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

func WReturnOrder(s s.OrderStorageInterface) error {
	id, err := inputs.InputOrderID()
	if err != nil {
		return err
	}

	return orders.ReturnOrder(s, id)
}

func WPlaceOrder(s s.OrderStorageInterface) error {
	uintdata, err := inputs.InputOrderIDs()
	if err != nil {
		return err
	}

	return orders.PlaceOrder(s, uintdata)
}

func WListOrders(s s.OrderStorageInterface) error {
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

func WRefundOrder(rS r.ReturnStorageInterface, oS s.OrderStorageInterface) error {
	orderId, userId, err := inputs.InputOrderAndUserID()
	if err != nil {
		return err
	}

	return returns.RefundOrder(rS, oS, orderId, userId)
}

func WListReturns(rs r.ReturnStorageInterface) error {
	limit, page, err := inputs.InputReturnsPagination()
	if err != nil {
		return err
	}

	return returns.ListReturns(rs, limit, page)
}
