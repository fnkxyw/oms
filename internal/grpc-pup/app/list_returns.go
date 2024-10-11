package pup_service

import (
	"context"
	desc "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) ListReturns(ctx context.Context, req *desc.ListReturnsRequest) (*desc.ListReturnsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	list, err := i.storage.ListReturns(ctx, int(req.Limit), int(req.Page))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := &desc.ListReturnsResponse{}
	for _, order := range list {
		response.Returns = append(response.Returns, &desc.Return{
			OrderId: uint32(order.ID),
			UserId:  uint32(order.UserID),
		})
	}
	return response, nil
}
