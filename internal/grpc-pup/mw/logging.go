package mw

import (
	"context"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
)

func Logging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Printf("[interceptor.Logging] method: %s; metadata: %v", info.FullMethod, md)
	}

	rewReq, _ := protojson.Marshal((req).(proto.Message))
	log.Printf("[interceptor.Logging] method: %s; request: %s", info.FullMethod, string(rewReq))

	res, err := handler(ctx, req)
	if err != nil {
		log.Printf("[interceptor.Logging] method: %s; error: %s", info.FullMethod, err.Error())
		metrics.IncBadRespByHandler(info.FullMethod)
		return
	}

	respReq, _ := protojson.Marshal((res).(proto.Message))
	log.Printf("[interceptor.Logging] method: %s; response: %s", info.FullMethod, string(respReq))
	metrics.IncOkRespByHandler(info.FullMethod)
	return res, nil
}
