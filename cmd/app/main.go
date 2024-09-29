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
	//oS := orderStorage.NewOrderStorage()
	const psqlDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	oS := newStorageFacade(pool)
	//oS.AddToStorage(ctx, &models.Order{ID: 3, UserID: 1, State: models.AcceptState})

	//err = signals.SignalSearch(*oS)
	//if err != nil {
	//	return
	//}

	err = c.Run(ctx, oS)
	if err != nil {
		return
	}
	//oS.WriteToJSON()
}

func newStorageFacade(pool *pgxpool.Pool) storage.Facade {
	txManager := postgres.NewTxManager(pool)

	pgRepository := postgres.NewPgRepository(txManager)

	return storage.NewStorageFacade(txManager, pgRepository)
}
