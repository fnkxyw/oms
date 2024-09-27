package postgres

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"time"
)

type PgRepository struct {
	pool *pgxpool.Pool
}

func NewPgRepositrory(pool *pgxpool.Pool) *PgRepository {
	return &PgRepository{pool: pool}
}

//AddToStorage(order *models.Order)
//IsConsist(id uint) bool
//DeleteFromStorage(id uint)
//GetItem(id uint) (*models.Order, bool)
//GetIDs() []uint

func (r *PgRepository) AddToStorage(ctx context.Context, order *models.Order) {
	_, err := r.pool.Exec(ctx, `INSERT INTO orders (id, user_id, state,accept_time,keep_until_date,place_date, weight,price) VALUES ($1, $2, $3,$4,$5,$6,$7,$8)`, order.ID, order.UserID, order.State, order.AcceptTime, order.KeepUntilDate, time.Time{}, order.Weight, order.Price)
	if err != nil {
		panic(err)
	}
}

func (r *PgRepository) IsConsist(ctx context.Context, id uint) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1)`
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (r *PgRepository) DeleteFromStorage(ctx context.Context, id uint) {
	_, err := r.pool.Exec(ctx, `UPDATE orders SET state =  $1 WHERE id = $2`, models.SoftDelete, id)
	if err != nil {
		panic(err)
	}
}

func (r *PgRepository) GetItem(ctx context.Context, id uint) (*models.Order, bool) {
	order := models.Order{}
	err := pgxscan.Get(ctx, r.pool, &order, `SELECT id, user_id, state,accept_time,keep_until_date,place_date,weight,price FROM orders WHERE id = $1`, id)
	if err != nil {
		return nil, false
	}
	return &order, true
}

func (r *PgRepository) GetIDs(ctx context.Context) []uint {
	ids := make([]uint, 0)
	err := pgxscan.Select(ctx, r.pool, &ids, `SELECT id FROM orders`)
	if err != nil {
		return nil
	}
	return ids
}

func (r *PgRepository) UpdateBeforePlace(ctx context.Context, id uint, state models.State, t time.Time) error {
	_, err := r.pool.Exec(ctx, `UPDATE orders SET state = $1, place_date = $2 WHERE id = $3`, state, t, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PgRepository) UpdateState(ctx context.Context, id uint, state models.State) error {
	_, err := r.pool.Exec(ctx, `UPDATE orders SET state = $1 WHERE id = $2`, state, id)
	if err != nil {
		return err
	}
	return nil
}
