package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	c "gitlab.ozon.dev/akugnerevich/homework.git/internal/cli"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	"log"
)

func main() {
	const psqlDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	const psqlDSNr = "postgres://replicator:replicator_password@localhost:5434/postgres"
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	poolRep, err := pgxpool.New(ctx, psqlDSNr)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()
	defer poolRep.Close()
	oS := storage.NewStorageFacade(pool)

	err = c.Run(ctx, oS)
	if err != nil {
		return
	}

}
