package pup_service

import (
	"context"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders"
	desc "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *Implementation) ListOrders(ctx context.Context, req *desc.ListOrdersRequest) (*desc.ListOrdersResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	list, err := i.storage.ListOrders(ctx, uint(req.UserId), req.InPup)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	count := int(req.Count)
	orders.SortOrders(list)
	if count < 1 {
		count = 1
	} else if count > len(list) {
		count = len(list)

	}
	if !req.InPup {
		list = list[:count]
	}
	response := &desc.ListOrdersResponse{}
	for _, order := range list {
		response.Orders = append(response.Orders, &desc.OrderFromList{
			OrderId:       uint32(order.ID),
			UserId:        uint32(order.UserID),
			State:         string(order.State),
			Price:         int32(order.Price),
			KeepUntilDate: timestamppb.New(order.KeepUntilDate),
		})
	}
	return response, nil
}
