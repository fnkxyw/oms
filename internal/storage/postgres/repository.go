package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"time"
)

var (
	ErrorNotFoundContact = errors.New("contact not found")
)

type PgRepository struct {
	txManager TransactionManager
}

func NewPgRepository(txManager TransactionManager) *PgRepository {
	return &PgRepository{
		txManager: txManager,
	}
}

func (r *PgRepository) AddToStorage(ctx context.Context, order *models.Order) error {
	tx := r.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, `INSERT INTO orders (id, user_id, state, accept_time, keep_until_date, place_date, weight, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.ID, order.UserID, order.State, order.AcceptTime, order.KeepUntilDate, order.PlaceDate, order.Weight, order.Price)
	return err
}

func (r *PgRepository) IsConsist(ctx context.Context, id uint) bool {
	tx := r.txManager.GetQueryEngine(ctx)
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1)`
	err := tx.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (r *PgRepository) DeleteFromStorage(ctx context.Context, id uint) error {
	tx := r.txManager.GetQueryEngine(ctx)
	result, err := tx.Exec(ctx, `UPDATE orders SET state = $1 WHERE id = $2`, models.SoftDelete, id)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrorNotFoundContact
	}

	return err
}

func (r *PgRepository) GetItem(ctx context.Context, id uint) (*models.Order, bool) {
	tx := r.txManager.GetQueryEngine(ctx)
	order := models.Order{}
	err := tx.QueryRow(ctx, `SELECT id, user_id, state, accept_time, keep_until_date, place_date, weight, price FROM orders WHERE id = $1`, id).Scan(
		&order.ID, &order.UserID, &order.State, &order.AcceptTime, &order.KeepUntilDate, &order.PlaceDate, &order.Weight, &order.Price,
	)
	if err != nil {
		return nil, false
	}
	return &order, true
}

func (r *PgRepository) UpdateBeforePlace(ctx context.Context, id uint, state models.State, t time.Time) error {
	tx := r.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, `UPDATE orders SET state = $1, place_date = $2 WHERE id = $3`, state, t, id)
	return err
}

func (r *PgRepository) UpdateState(ctx context.Context, id uint, state models.State) error {
	tx := r.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, `UPDATE orders SET state = $1 WHERE id = $2`, state, id)
	return err
}

func (r *PgRepository) GetOrders(ctx context.Context, userId uint, inPuP bool) ([]models.Order, error) {
	tx := r.txManager.GetQueryEngine(ctx)
	rows, err := tx.Query(ctx, `SELECT * FROM orders WHERE user_id = $1`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Order])
	if err != nil {
		return nil, err
	}
	var filtered []models.Order

	for _, v := range orders {
		if !inPuP || v.State == models.AcceptState || v.State == models.RefundedState {
			filtered = append(filtered, v)
		}
	}
	return orders, nil
}

func (r *PgRepository) GetReturns(ctx context.Context, page, limit int) ([]models.Order, error) {
	tx := r.txManager.GetQueryEngine(ctx)
	rows, err := tx.Query(ctx, `SELECT * FROM orders WHERE state = $1 LIMIT $2 OFFSET $3`, models.RefundedState, limit, page*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	returns, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Order])
	if err != nil {
		return nil, err
	}

	return returns, nil
}
