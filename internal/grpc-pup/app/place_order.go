package pup_service

import (
	"context"
	desc "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) PlaceOrderV1(ctx context.Context, req *desc.PlaceOrderV1Request) (*desc.PlaceOrderV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := i.storage.PlaceOrder(ctx, req.OrderIds)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &desc.PlaceOrderV1Response{}, nil
}
