package pup_service

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/metrics"
	desc "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) ReturnOrderV1(ctx context.Context, req *desc.ReturnOrderV1Request) (*desc.ReturnOrderV1Response, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReturnOrderV1")
	defer span.Finish()
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := i.storage.ReturnOrder(ctx, uint(req.OrderId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	metrics.IncOrderTotalOperations("return")

	return nil, nil
}
