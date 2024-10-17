package pup_service

import (
	"context"
	desc "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) ListReturnsV1(ctx context.Context, req *desc.ListReturnsV1Request) (*desc.ListReturnsV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	list, err := i.storage.ListReturns(ctx, int(req.Limit), int(req.Page))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := &desc.ListReturnsV1Response{}
	for _, order := range list {
		response.Returns = append(response.Returns, &desc.ReturnV1{
			OrderId: uint32(order.ID),
			UserId:  uint32(order.UserID),
		})
	}
	return response, nil
}
