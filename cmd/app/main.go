package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	c "gitlab.ozon.dev/akugnerevich/homework.git/internal/cli"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/postgres"
	"log"
)

func main() {
	//oS := orderStorage.NewOrderStorage()
	const psqlDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	oS := postgres.NewPgRepositrory(pool)
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
