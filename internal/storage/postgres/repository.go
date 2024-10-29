package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "PgRepository.AddToStorage")
	defer span.Finish()

	tx := r.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, `INSERT INTO orders (id, user_id, state, accept_time, keep_until_date, place_date, weight, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.ID, order.UserID, order.State, order.AcceptTime, order.KeepUntilDate, order.PlaceDate, order.Weight, order.Price)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			span.LogKV("error", e.ErrIsConsist)
			return e.ErrIsConsist
		}
		span.LogKV("error", err.Error())
		return err
	}

	span.LogKV("info", "order added to storage", "order_id", order.ID)
	return nil
}

func (r *PgRepository) DeleteFromStorage(ctx context.Context, id uint) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PgRepository.DeleteFromStorage")
	defer span.Finish()

	tx := r.txManager.GetQueryEngine(ctx)
	result, err := tx.Exec(ctx, `UPDATE orders SET state = $1 WHERE id = $2`, models.SoftDelete, id)
	if err != nil {
		span.LogKV("error", err.Error())
		return err
	}
	if result.RowsAffected() == 0 {
		span.LogKV("error", ErrorNotFoundOrder)
		return ErrorNotFoundOrder
	}

	span.LogKV("info", "order deleted", "order_id", id)
	return nil
}

func (r *PgRepository) GetItem(ctx context.Context, id uint) (*models.Order, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PgRepository.GetItem")
	defer span.Finish()

	tx := r.txManager.GetQueryEngine(ctx)
	const query = `SELECT ` + orderFields + ` FROM orders WHERE id = $1 FOR UPDATE`
	rows, err := tx.Query(ctx, query, id)
	if err != nil {
		span.LogKV("error", err.Error())
		return nil, false
	}

	order, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[models.Order])
	if err != nil {
		span.LogKV("error", err.Error())
		return nil, false
	}

	span.LogKV("info", "order retrieved", "order_id", id)
	return &order, true
}

func (r *PgRepository) UpdateBeforePlace(ctx context.Context, ids []uint32, t time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PgRepository.UpdateBeforePlace")
	defer span.Finish()

	tx := r.txManager.GetQueryEngine(ctx)
	query := `UPDATE orders SET state = $1, place_date = $2 WHERE id = ANY($3)`
	_, err := tx.Exec(ctx, query, models.PlaceState, t, ids)
	if err != nil {
		span.LogKV("error", err.Error())
	}
	return err
}

func (r *PgRepository) UpdateState(ctx context.Context, id uint, state models.State) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PgRepository.UpdateState")
	defer span.Finish()

	tx := r.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, `UPDATE orders SET state = $1 WHERE id = $2`, state, id)
	if err != nil {
		span.LogKV("error", err.Error())
	}
	return err
}

func (r *PgRepository) GetUserOrders(ctx context.Context, userId uint, inPuP bool) ([]models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PgRepository.GetUserOrders")
	defer span.Finish()

	tx := r.txManager.GetQueryEngine(ctx)
	const query = `SELECT ` + orderFields + ` FROM orders WHERE user_id = $1`
	rows, err := tx.Query(ctx, query, userId)
	if err != nil {
		span.LogKV("error", err.Error())
		return nil, err
	}

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Order])
	if err != nil {
		span.LogKV("error", err.Error())
		return nil, err
	}

	var filtered []models.Order
	for _, v := range orders {
		if !inPuP || v.State == models.AcceptState || v.State == models.RefundedState {
			filtered = append(filtered, v)
		}
	}

	span.LogKV("info", "user orders retrieved", "user_id", userId)
	return filtered, nil
}

func (r *PgRepository) GetReturns(ctx context.Context, page, limit int) ([]models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PgRepository.GetReturns")
	defer span.Finish()

	tx := r.txManager.GetQueryEngine(ctx)
	const query = `SELECT ` + orderFields + ` FROM orders WHERE state = $1 LIMIT $2 OFFSET $3`
	rows, err := tx.Query(ctx, query, models.RefundedState, limit, page*limit)
	if err != nil {
		span.LogKV("error", err.Error())
		return nil, err
	}

	returns, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Order])
	if err != nil {
		span.LogKV("error", err.Error())
		return nil, err
	}

	span.LogKV("info", "returns retrieved", "page", page, "limit", limit)
	return returns, nil
}

func (r *PgRepository) GetItems(ctx context.Context, ids []uint32) ([]models.Order, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PgRepository.GetItems")
	defer span.Finish()

	tx := r.txManager.GetQueryEngine(ctx)
	const query = `SELECT ` + orderFields + ` FROM orders WHERE id = ANY($1) FOR UPDATE`
	rows, err := tx.Query(ctx, query, ids)
	if err != nil {
		span.LogKV("error", err.Error())
		return nil, false
	}

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.Order])
	if err != nil {
		span.LogKV("error", err.Error())
		return nil, false
	}

	if len(orders) != len(ids) {
		span.LogKV("error", "not all requested IDs found")
		return orders, false
	}

	span.LogKV("info", "items retrieved", "count", len(orders))
	return orders, true
}
