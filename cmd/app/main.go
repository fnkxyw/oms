package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	c "gitlab.ozon.dev/akugnerevich/homework.git/internal/cli"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/postgres"
	"log"
)

func main() {
	const psqlDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	oS := newStorageFacade(pool)

	//err = signals.SignalSearch(*oS)
	//if err != nil {
	//	return
	//}

	err = c.Run(ctx, oS)
	if err != nil {
		return
	}

}

func newStorageFacade(pool *pgxpool.Pool) storage.Facade {
	txManager := postgres.NewTxManager(pool)

	pgRepository := postgres.NewPgRepository(txManager)

	return storage.NewStorageFacade(txManager, pgRepository)
}
