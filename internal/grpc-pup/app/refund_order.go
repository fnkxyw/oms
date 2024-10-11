package pup_service

import (
	"context"
	desc "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) RefundOrder(ctx context.Context, req *desc.RefundOrderRequest) (*desc.RefundOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := i.storage.RefundOrder(ctx, uint(req.OrderId), uint(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}
