package main

import (
	"context"
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
	"log"
	"net"
	"net/http"
	"os"
)

const (
	psqlDBN   = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	adminHost = "localhost:7003"
	grpcHost  = "localhost:7002"
	httpHost  = "localhost:7001"
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

	go func() {
		adminServer := chi.NewMux()

		adminServer.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			b, err := os.ReadFile("./pkg/PuP-service/v1/pup_service.swagger.json")
			if err != nil {
				http.Error(w, "Could not read swagger file", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		})

		adminServer.Get("/swagger-ui", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./swagger_ui.html")
		})

		if err := http.ListenAndServe(adminHost, adminServer); err != nil {
			log.Fatalf("failed to listen and serve admin server: %v", err)
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
