package int_Postgres

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage"
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
			wait.ForAll(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2),
				wait.ForListeningPort(nat.Port("5432/tcp")),
			),
		),
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

func TestAcceptOrder(t *testing.T) {
	storageFacade := storage.NewStorageFacade(pool)

	order := &models.Order{
		ID:            4,
		UserID:        1,
		State:         models.NewState,
		AcceptTime:    time.Now().Unix(),
		KeepUntilDate: time.Now().Add(24 * time.Hour),
		Weight:        1.0,
		Price:         15.0,
	}

	err := storageFacade.AcceptOrder(context.Background(), order)
	assert.NoError(t, err, "Expected no error when accepting order")
	acceptedOrder, exists := storageFacade.GetItem(context.Background(), order.ID)
	assert.True(t, exists)
	assert.Equal(t, models.AcceptState, acceptedOrder.State)

}

func TestPlaceOrder(t *testing.T) {

	storageFacade := storage.NewStorageFacade(pool)

	orderIDs := []uint{1, 2}

	err := storageFacade.PlaceOrder(context.Background(), orderIDs)
	assert.NoError(t, err)

	for _, id := range orderIDs {
		order, exists := storageFacade.GetItem(context.Background(), id)
		assert.True(t, exists)
		assert.Equal(t, models.PlaceState, order.State)
	}
}

func TestReturnOrder(t *testing.T) {

	storageFacade := storage.NewStorageFacade(pool)

	orderID := uint(3)

	err := storageFacade.ReturnOrder(context.Background(), orderID)
	assert.Error(t, err)

	order, exists := storageFacade.GetItem(context.Background(), orderID)
	assert.True(t, exists)
	assert.Equal(t, models.PlaceState, order.State)
}

func TestRefundOrder(t *testing.T) {

	storageFacade := storage.NewStorageFacade(pool)

	orderID := uint(3)
	userID := uint(2)

	err := storageFacade.RefundOrder(context.Background(), orderID, userID)
	assert.NoError(t, err)

	order, exists := storageFacade.GetItem(context.Background(), orderID)
	assert.True(t, exists)
	assert.Equal(t, models.RefundedState, order.State)
}

func TestListOrders(t *testing.T) {

	storageFacade := storage.NewStorageFacade(pool)

	userID := uint(1)

	orders, err := storageFacade.ListOrders(context.Background(), userID, true)
	assert.NoError(t, err)
	assert.Greater(t, len(orders), 0)
}

func TestListReturns(t *testing.T) {

	storageFacade := storage.NewStorageFacade(pool)

	returns, err := storageFacade.ListReturns(context.Background(), 10, 1)
	assert.NoError(t, err)
	assert.Greater(t, len(returns), 0)
}

func TestTeardown(t *testing.T) {
	pool.Close()
	if err := pgContainer.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}
