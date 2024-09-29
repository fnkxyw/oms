package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type QueryEngine interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

type TransactionManager interface {
	GetQueryEngine(ctx context.Context) QueryEngine
	RunReadUncommitted(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunSerializable(ctx context.Context, fn func(ctxTx context.Context) error) error
}
