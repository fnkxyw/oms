package main

import (
	"context"
	_ "embed"
	kafka "gitlab.ozon.dev/akugnerevich/homework.git/internal/kafka/sync_producer"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	pup_service "gitlab.ozon.dev/akugnerevich/homework.git/internal/grpc-pup/app"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/grpc-pup/mw"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	desc "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	psqlDBN   = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	adminHost = "localhost:7003"
	grpcHost  = "localhost:7002"
	httpHost  = "localhost:7001"
	kafkaHost = "localhost:9092"
)

//go:embed swagger/pup_service.swagger.json
var swaggerJSON []byte

//go:embed swagger/swagger_ui.html
var swaggerUI []byte

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDBN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	producer, err := kafka.NewSyncProducer([]string{kafkaHost}, "pup-topic")
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}
	defer producer.Close()
	storageFacade := storage.NewStorageFacade(pool, producer)
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

	go func() {
		adminServer := chi.NewMux()

		adminServer.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write(swaggerJSON)
		})

		adminServer.Get("/swagger-ui", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write(swaggerUI)
		})

		if err := http.ListenAndServe(adminHost, adminServer); err != nil {
			log.Fatalf("failed to listen and serve admin server: %v", err)
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
