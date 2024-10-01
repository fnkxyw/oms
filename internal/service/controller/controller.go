package controller

import (
	"context"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/controller/inputs"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders/packing"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/returns"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
)

func WAcceptOrder(ctx context.Context, s storage.Facade) error {
	order, packageType, needWrapping, err := inputs.CollectOrderInput()
	if err != nil {
		return err
	}

	err = packing.Packing(order, packageType, needWrapping)
	if err != nil {
		return err
	}

	err = orders.AcceptOrder(ctx, s, order)
	if err != nil {
		return err
	}

	return nil
}

func WReturnOrder(ctx context.Context, s storage.Facade) error {
	id, err := inputs.InputOrderID()
	if err != nil {
		return err
	}

	return orders.ReturnOrder(ctx, s, id)
}

func WPlaceOrder(ctx context.Context, s storage.Facade) error {
	uintdata, err := inputs.InputOrderIDs()
	if err != nil {
		return err
	}

	return orders.PlaceOrder(ctx, s, uintdata)
}

func WListOrders(ctx context.Context, s storage.Facade) error {
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
		return orders.ListOrders(ctx, s, id, 0, true)
	case 2:
		n, err := inputs.InputN()
		if err != nil {
			return err
		}
		return orders.ListOrders(ctx, s, id, n, false)
	}

	return nil
}

func WRefundOrder(ctx context.Context, oS storage.Facade) error {
	orderId, userId, err := inputs.InputOrderAndUserID()
	if err != nil {
		return err
	}

	return returns.RefundOrder(ctx, oS, orderId, userId)
}

func WListReturns(ctx context.Context, oS storage.Facade) error {
	limit, page, err := inputs.InputReturnsPagination()
	if err != nil {
		return err
	}

	return returns.ListReturns(ctx, oS, limit, page)
}
