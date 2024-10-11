package pup_service

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	pup_service "gitlab.ozon.dev/akugnerevich/homework.git/internal/grpc-pup/app"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	"log"
)

const (
	psqlDBN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
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
}
