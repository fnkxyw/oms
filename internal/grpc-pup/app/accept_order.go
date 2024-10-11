package pup_service

import (
	"context"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders/packing"
	desc "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (i *Implementation) AcceptOrder(ctx context.Context, req *desc.AcceptOrderRequest) (*desc.AcceptOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	order := &models.Order{
		ID:            uint(req.OrderId),
		UserID:        uint(req.UserId),
		State:         "",
		AcceptTime:    time.Now().Unix(),
		KeepUntilDate: req.KeepUntilDate.AsTime(),
		PlaceDate:     time.Time{},
		Weight:        int(req.Weight),
		Price:         int(req.Price),
	}

	err := packing.Packing(order, req.PackageType, req.NeedWrapping)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = i.storage.AcceptOrder(ctx, order)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &desc.AcceptOrderResponse{}, nil
}
