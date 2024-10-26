package pup_service

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/metrics"
	desc "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *Implementation) ListOrdersV1(ctx context.Context, req *desc.ListOrdersV1Request) (*desc.ListOrdersV1Response, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ListOrdersV1")
	defer span.Finish()
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	list, err := i.storage.ListOrders(ctx, uint(req.UserId), req.InPup, int(req.Count))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &desc.ListOrdersV1Response{}
	for _, order := range list {
		response.Orders = append(response.Orders, &desc.OrderFromListV1{
			OrderId:       uint32(order.ID),
			UserId:        uint32(order.UserID),
			State:         string(order.State),
			Price:         int32(order.Price),
			KeepUntilDate: timestamppb.New(order.KeepUntilDate),
		})
	}
	metrics.IncOrderTotalOperations("list_orders")

	return response, nil
}
