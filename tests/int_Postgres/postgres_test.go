package int_Postgres

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/postgres"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	pl "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbName     = "test_db"
	dbUser     = "test_user"
	dbPassword = "test_password"
)

var (
	pgContainer testcontainers.Container
	pool        *pgxpool.Pool
	ctx         = context.Background()
)

func TestSetup(t *testing.T) {
	var err error

	pgContainer, err = pl.Run(ctx, "docker.io/postgres:15-alpine",
		pl.WithDatabase(dbName),
		pl.WithUsername(dbUser),
		pl.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to open container: %s", err)
	}

	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to map port: %s", err)
	}

	dbURI := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPassword, mappedPort.Port(), dbName)

	pool, err = pgxpool.New(ctx, dbURI)
	if err != nil {
		t.Fatalf("failed to connect to database: %s", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	time.Sleep(1 * time.Second)

	// Run migrations
	if err := goose.Up(db, "../../migrations"); err != nil {
		t.Fatalf("failed to run migrations: %s", err)
	}

	if err := seedTestData(ctx, pool); err != nil {
		t.Fatalf("failed to seed test data: %s", err)
	}
}

func seedTestData(ctx context.Context, pool *pgxpool.Pool) error {
	orders := []models.Order{
		{
			ID:            1,
			UserID:        1,
			State:         models.AcceptState,
			AcceptTime:    time.Now().Unix(),
			KeepUntilDate: time.Now().Add(24 * time.Hour),
			PlaceDate:     time.Now(),
			Weight:        1.0,
			Price:         10.0,
		},
		{
			ID:            2,
			UserID:        1,
			State:         models.NewState,
			AcceptTime:    time.Now().Unix(),
			KeepUntilDate: time.Now().Add(24 * time.Hour),
			PlaceDate:     time.Now(),
			Weight:        2.0,
			Price:         20.0,
		},
		{
			ID:            3,
			UserID:        2,
			State:         models.PlaceState,
			AcceptTime:    time.Now().Unix(),
			KeepUntilDate: time.Now().Add(48 * time.Hour),
			PlaceDate:     time.Now(),
			Weight:        2,
			Price:         30.0,
		},
	}

	for _, order := range orders {
		_, err := pool.Exec(ctx, `INSERT INTO orders (id, user_id, state, accept_time, keep_until_date, place_date, weight, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			order.ID, order.UserID, order.State, order.AcceptTime, order.KeepUntilDate, order.PlaceDate, order.Weight, order.Price)
		if err != nil {
			return fmt.Errorf("failed to insert order %v: %w", order, err)
		}
	}

	return nil
}

func TestAddToStorage(t *testing.T) {
	repo := postgres.NewPgRepository(pool)

	order := &models.Order{
		ID:            4,
		UserID:        1,
		State:         models.NewState,
		AcceptTime:    time.Now().Unix(),
		KeepUntilDate: time.Now().Add(24 * time.Hour),
		PlaceDate:     time.Now(),
		Weight:        3.0,
		Price:         40.0,
	}

	repo.AddToStorage(ctx, order)

	item, found := repo.GetItem(ctx, 4)
	if !found || item.ID != 4 {
		t.Fatalf("expected to find order with ID 4, got %v", item)
	}
}

func TestIsConsist(t *testing.T) {
	repo := postgres.NewPgRepository(pool)

	exists := repo.IsConsist(ctx, 1)
	if !exists {
		t.Fatal("expected order with ID 1 to exist")
	}

	exists = repo.IsConsist(ctx, 100)
	if exists {
		t.Fatal("expected order with ID 100 to not exist")
	}
}

func TestDeleteFromStorage(t *testing.T) {
	repo := postgres.NewPgRepository(pool)

	repo.DeleteFromStorage(ctx, 1)

	deletedOrder, found := repo.GetItem(ctx, 1)
	if !found || deletedOrder.State != models.SoftDelete {
		t.Fatalf("expected order with ID 1 to be marked as deleted, but found: %v", deletedOrder)
	}
}

func TestGetItem(t *testing.T) {
	repo := postgres.NewPgRepository(pool)

	order, found := repo.GetItem(ctx, 2)
	if !found || order.ID != 2 {
		t.Fatalf("expected to find order with ID 2, got %v", order)
	}
}

func TestUpdateState(t *testing.T) {
	repo := postgres.NewPgRepository(pool)

	err := repo.UpdateState(ctx, 2, models.AcceptState)
	if err != nil {
		t.Fatalf("failed to update state: %s", err)
	}

	updatedOrder, _ := repo.GetItem(ctx, 2)
	if updatedOrder.State != models.AcceptState {
		t.Fatalf("expected state to be %s, got %s", models.AcceptState, updatedOrder.State)
	}
}

func TestUpdateBeforePlace(t *testing.T) {
	repo := postgres.NewPgRepository(pool)

	err := repo.UpdateBeforePlace(ctx, 2, models.PlaceState, time.Now())
	if err != nil {
		t.Fatalf("failed to update order before placing: %s", err)
	}

	updatedOrder, _ := repo.GetItem(ctx, 2)
	if updatedOrder.State != models.PlaceState {
		t.Fatalf("expected state to be %s, got %s", models.PlaceState, updatedOrder.State)
	}
}

func TestGetByUserId(t *testing.T) {
	repo := postgres.NewPgRepository(pool)

	orders, err := repo.GetByUserId(ctx, 1)
	if err != nil || len(orders) == 0 {
		t.Fatalf("expected to find orders for user ID 1, got error: %v", err)
	}
}

func TestGetReturns(t *testing.T) {
	repo := postgres.NewPgRepository(pool)

	resultOrders, err := repo.GetReturns(ctx, models.RefundedState)
	if err != nil {
		t.Fatalf("failed to get orders: %s", err)
	}

	expectedCount := 1
	if len(resultOrders) != expectedCount {
		t.Fatalf("expected %d orders, got %d", expectedCount, len(resultOrders))
	}

	for _, order := range resultOrders {
		if order.State != models.NewState {
			t.Fatalf("expected order state to be %s, got %s", models.NewState, order.State)
		}
	}
}

func TestTeardown(t *testing.T) {
	pool.Close()
	if err := pgContainer.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}
