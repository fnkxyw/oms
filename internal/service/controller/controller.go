package controller

import (
	"context"
	"github.com/containerd/containerd/protobuf"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/controller/inputs"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/wpool"
	pup_service "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
)

func logResponse(resp proto.Message, err error) {
	data, marshalErr := protojson.Marshal(resp)
	if marshalErr != nil {
		log.Printf("failed to marshal response: %v; original error: %v\n", marshalErr, err)
		return
	}

	log.Printf("resp: %s; err: %v\n", string(data), err)
}

func WAcceptOrder(ctx context.Context, pup pup_service.PupServiceClient, wp *wpool.WorkerPool) error {
	order, PackageType, NeedWrapping, err := inputs.CollectOrderInput()
	if err != nil {
		return err
	}

	req := &pup_service.AcceptOrderRequest{
		OrderId:       uint32(order.ID),
		UserId:        uint32(order.UserID),
		KeepUntilDate: protobuf.ToTimestamp(order.KeepUntilDate),
		Weight:        int32(order.Weight),
		Price:         int32(order.Price),
		PackageType:   PackageType,
		NeedWrapping:  NeedWrapping,
	}

	wp.AddJob(ctx, func() {
		resp, err := pup.AcceptOrder(ctx, req)
		logResponse(resp, err)
	}, "Adding and Packaging Order", 6)

	return nil
}

func WReturnOrder(ctx context.Context, pup pup_service.PupServiceClient, wp *wpool.WorkerPool) error {
	id, err := inputs.InputOrderID()
	if err != nil {
		return err
	}
	req := &pup_service.ReturnOrderRequest{OrderId: uint32(id)}
	wp.AddJob(ctx, func() {
		resp, err := pup.ReturnOrder(ctx, req)
		logResponse(resp, err)
	}, "Returning Order", 3)

	return nil
}

func WPlaceOrder(ctx context.Context, pup pup_service.PupServiceClient, wp *wpool.WorkerPool) error {
	uintdata, err := inputs.InputOrderIDs()
	if err != nil {
		return err
	}

	req := &pup_service.PlaceOrderRequest{OrderId: uintdata}

	wp.AddJob(ctx, func() {
		resp, err := pup.PlaceOrder(ctx, req)
		logResponse(resp, err)
	}, "Placing Order", 5)

	return nil
}

func WListOrders(ctx context.Context, pup pup_service.PupServiceClient, wp *wpool.WorkerPool) error {
	id, err := inputs.InputUserID()
	if err != nil {
		return err
	}
	temp, err := inputs.InputListChoice()
	if err != nil {
		return err
	}

	wp.AddJob(ctx, func() {
		switch temp {
		case 1:
			req := &pup_service.ListOrdersRequest{
				UserId: uint32(id),
				Count:  0,
				InPup:  true,
			}
			resp, err := pup.ListOrders(ctx, req)
			logResponse(resp, err)
		case 2:
			n, err := inputs.InputN()
			if err != nil {
				log.Printf("listOrders input failed: %v", err)
			}
			req := &pup_service.ListOrdersRequest{
				UserId: uint32(id),
				Count:  int32(n),
				InPup:  false,
			}
			resp, err := pup.ListOrders(ctx, req)
			logResponse(resp, err)
		default:
			log.Printf("invalid input")
			return
		}

	}, "Listing Orders", 2)

	return nil
}

func WRefundOrder(ctx context.Context, pup pup_service.PupServiceClient, wp *wpool.WorkerPool) error {
	orderId, userId, err := inputs.InputOrderAndUserID()
	if err != nil {
		return err
	}

	req := &pup_service.RefundOrderRequest{
		OrderId: uint32(orderId),
		UserId:  uint32(userId),
	}

	wp.AddJob(ctx, func() {

		resp, err := pup.RefundOrder(ctx, req)
		logResponse(resp, err)
	}, "Refunding Order", 4)

	return nil
}

func WListReturns(ctx context.Context, pup pup_service.PupServiceClient, wp *wpool.WorkerPool) error {

	limit, page, err := inputs.InputReturnsPagination()
	if err != nil {
		return err
	}

	req := &pup_service.ListReturnsRequest{
		Limit: int32(limit),
		Page:  int32(page),
	}

	wp.AddJob(ctx, func() {

		resp, err := pup.ListReturns(ctx, req)
		logResponse(resp, err)
	}, "Listing Returns", 1)

	return nil
}

func WChangeNumOfWorkers(ctx context.Context, wp *wpool.WorkerPool) error {
	n, err := inputs.InputNumOfWorkers(wp)
	if err != nil {
		return err
	}

	wp.AddJob(ctx, func() {
		if err := wp.ChangeNumOfWorkers(n); err != nil {
			log.Printf("change num of workers error: %v", err)
		}
		wp.PrintWorkers()
	}, "Changing Number of Workers", 3)

	return nil
}
