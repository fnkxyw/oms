package main

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/postgres"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	const psqlDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	rep := postgres.NewPgRepositrory(pool)

	for i := 0; i < 100000; i++ {
		order := models.Order{
			ID:            uint(i + 1),
			UserID:        uint(gofakeit.Number(1, 10)),
			State:         generateState(),
			AcceptTime:    int64(gofakeit.Number(int(time.Now().Unix()), int(time.Now().Unix()+1000))),
			KeepUntilDate: gofakeit.Date(),
			PlaceDate:     gofakeit.Date(),
			Weight:        gofakeit.Number(1, 100),
			Price:         gofakeit.Number(10, 10000),
		}

		rep.AddToStorage(ctx, &order)
	}

}
