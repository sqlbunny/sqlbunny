package boil

import (
	"context"
	"database/sql"
)

// Executor can perform SQL queries.
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type key int

const dbKey key = 0

func WithDB(ctx context.Context, db Executor) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

func DBFromContext(ctx context.Context) Executor {
	db, ok := ctx.Value(dbKey).(Executor)
	if !ok {
		panic("No database in the context")
	}
	return db
}
