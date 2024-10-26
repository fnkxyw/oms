package pup_service

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/metrics"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders/packing"
	desc "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (i *Implementation) AcceptOrderV1(ctx context.Context, req *desc.AcceptOrderV1Request) (*desc.AcceptOrderV1Response, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AcceptOrderV1")
	defer span.Finish()
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	order := &models.Order{
		ID:            uint(req.OrderId),
		UserID:        uint(req.UserId),
		State:         "",
		AcceptTime:    time.Time{}.Unix(),
		KeepUntilDate: req.KeepUntilDate.AsTime(),
		PlaceDate:     time.Time{},
		Weight:        int(req.Weight),
		Price:         int(req.Price),
	}

	err := packing.Packing(order, packing.PackageType(req.PackageType), req.NeedWrapping)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = i.storage.AcceptOrder(ctx, order)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	metrics.IncOrderTotalOperations("accept")
	return &desc.AcceptOrderV1Response{}, nil
}
