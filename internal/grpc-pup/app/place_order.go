package pup_service

import (
	"context"
	desc "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) PlaceOrder(ctx context.Context, req *desc.PlaceOrderRequest) (*desc.PlaceOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := i.storage.PlaceOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &desc.PlaceOrderResponse{}, nil
}
