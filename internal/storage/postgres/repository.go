package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"time"
)

var (
	ErrorNotFoundOrder = errors.New("order not found")
)

const orderFields = "id, user_id, state, accept_time, keep_until_date, place_date, weight, price"

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

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return e.ErrIsConsist
			}
		}
		return err
	}

	return nil
}

func (r *PgRepository) DeleteFromStorage(ctx context.Context, id uint) error {
	tx := r.txManager.GetQueryEngine(ctx)
	result, err := tx.Exec(ctx, `UPDATE orders SET state = $1 WHERE id = $2`, models.SoftDelete, id)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrorNotFoundOrder
	}

	return err
}

func (r *PgRepository) GetItem(ctx context.Context, id uint) (*models.Order, bool) {
	tx := r.txManager.GetQueryEngine(ctx)
	const query = `SELECT ` + orderFields + ` FROM orders WHERE id = $1 FOR UPDATE`
	rows, err := tx.Query(ctx, query, id)
	if err != nil {
		return nil, false
	}

	order, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[models.Order])
	if err != nil {
		return nil, false
	}

	return &order, true
}

func (r *PgRepository) UpdateBeforePlace(ctx context.Context, ids []uint, t time.Time) error {
	tx := r.txManager.GetQueryEngine(ctx)

	query := `UPDATE orders	SET state = $1, place_date = $2	WHERE id = ANY($3)`
	_, err := tx.Exec(ctx, query, models.PlaceState, t, ids)
	return err
}

func (r *PgRepository) UpdateState(ctx context.Context, id uint, state models.State) error {
	tx := r.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, `UPDATE orders SET state = $1 WHERE id = $2`, state, id)
	return err
}

func (r *PgRepository) GetUserOrders(ctx context.Context, userId uint, inPuP bool) ([]models.Order, error) {
	tx := r.txManager.GetQueryEngine(ctx)
	const query = `SELECT ` + orderFields + ` FROM orders WHERE user_id = $1`
	rows, err := tx.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}

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
	const query = `SELECT ` + orderFields + ` FROM orders WHERE state = $1 LIMIT $2 OFFSET $3`

	rows, err := tx.Query(ctx, query, models.RefundedState, limit, page*limit)
	if err != nil {
		return nil, err
	}

	returns, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Order])
	if err != nil {
		return nil, err
	}

	return returns, nil
}

func (r *PgRepository) GetItems(ctx context.Context, ids []uint) ([]models.Order, bool) {
	tx := r.txManager.GetQueryEngine(ctx)
	const query = `SELECT  ` + orderFields + ` FROM orders WHERE id = ANY($1) FOR UPDATE`

	rows, err := tx.Query(ctx, query, ids)
	if err != nil {
		return nil, false
	}

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.Order])
	if err != nil {
		return nil, false
	}

	if len(orders) != len(ids) {
		return orders, false
	}

	return orders, true
}
