package main

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	const psqlDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	rep := storage.NewStorageFacade(pool)
	for i := 0; i < 200000; i++ {
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

		err := rep.PgRepo.AddToStorage(ctx, &order)
		if err != nil {
			log.Fatal(err)
		}
	}

}
