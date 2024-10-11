package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	pup_service "gitlab.ozon.dev/akugnerevich/homework.git/internal/grpc-pup/app"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/grpc-pup/mw"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	desc "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

const (
	psqlDBN  = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	grpcHost = "localhost:7001"
	httpHost = "localhost:7000"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDBN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	storageFacade := storage.NewStorageFacade(pool)
	pupService := pup_service.NewImplementation(storageFacade)

	lis, err := net.Listen("tcp", grpcHost)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mw.Logging),
	)
	reflection.Register(grpcServer)
	desc.RegisterPupServiceServer(grpcServer, pupService)

	mux := runtime.NewServeMux()
	err = desc.RegisterPupServiceHandlerFromEndpoint(ctx, mux, grpcHost, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		log.Fatalf("failed to register pup service handler: %v", err)
	}

	go func() {
		if err := http.ListenAndServe(httpHost, mux); err != nil {
			log.Fatalf("failed to listen and serve pup service handler: %v", err)
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
