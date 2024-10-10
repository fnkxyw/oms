package controller

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/controller/inputs"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders/packing"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/returns"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/wpool"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
)

func WAcceptOrder(ctx context.Context, s storage.Facade, wp *wpool.WorkerPool, errChan chan error) error {
	order, packageType, needWrapping, err := inputs.CollectOrderInput()
	if err != nil {
		return err
	}

	wp.AddJob(ctx, func() {
		if err := packing.Packing(order, packageType, needWrapping); err != nil {
			errChan <- err
			return
		}

		if err := orders.AcceptOrder(ctx, s, order); err != nil {
			errChan <- err
		}
	}, "Adding and Packaging Order", 6)

	return nil
}

func WReturnOrder(ctx context.Context, s storage.Facade, wp *wpool.WorkerPool, errChan chan error) error {
	id, err := inputs.InputOrderID()
	if err != nil {
		return err
	}

	wp.AddJob(ctx, func() {
		if err := orders.ReturnOrder(ctx, s, id); err != nil {
			errChan <- err
		}
	}, "Returning Order", 3)

	return nil
}

func WPlaceOrder(ctx context.Context, s storage.Facade, wp *wpool.WorkerPool, errChan chan error) error {
	uintdata, err := inputs.InputOrderIDs()
	if err != nil {
		return err
	}

	wp.AddJob(ctx, func() {

		if err := orders.PlaceOrder(ctx, s, uintdata); err != nil {
			errChan <- err

		}
	}, "Placing Order", 5)

	return nil
}

func WListOrders(ctx context.Context, s storage.Facade, wp *wpool.WorkerPool, errChan chan error) error {
	id, err := inputs.InputUserID()
	if err != nil {
		return err
	}

	temp, err := inputs.InputListChoice()
	if err != nil {
		return err
	}

	wp.AddJob(ctx, func() {
		var result error
		switch temp {
		case 1:
			result = orders.ListOrders(ctx, s, id, 0, true)
		case 2:
			n, err := inputs.InputN()
			if err != nil {
				errChan <- err
				return
			}
			result = orders.ListOrders(ctx, s, id, n, false)
		default:
			errChan <- fmt.Errorf("invalid choice")
			return
		}

		if result != nil {
			errChan <- result
		}
	}, "Listing Orders", 2)

	return nil
}

func WRefundOrder(ctx context.Context, oS storage.Facade, wp *wpool.WorkerPool, errChan chan error) error {
	orderId, userId, err := inputs.InputOrderAndUserID()
	if err != nil {
		return err
	}

	wp.AddJob(ctx, func() {
		if err := returns.RefundOrder(ctx, oS, orderId, userId); err != nil {
			errChan <- err
		}
	}, "Refunding Order", 4)

	return nil
}

func WListReturns(ctx context.Context, oS storage.Facade, wp *wpool.WorkerPool, errChan chan error) error {

	limit, page, err := inputs.InputReturnsPagination()
	if err != nil {
		return err
	}

	wp.AddJob(ctx, func() {
		if err := returns.ListReturns(ctx, oS, limit, page); err != nil {
			errChan <- err
		}
	}, "Listing Returns", 1)

	return nil
}

func WChangeNumOfWorkers(ctx context.Context, wp *wpool.WorkerPool, errChan chan error) error {
	n, err := inputs.InputNumOfWorkers(wp)
	if err != nil {
		return err
	}

	wp.AddJob(ctx, func() {
		if err := wp.ChangeNumOfWorkers(n); err != nil {
			errChan <- err
		}
		wp.PrintWorkers()
	}, "Changing Number of Workers", 3)

	return nil
}
