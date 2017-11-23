package boil

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// Executor can perform SQL queries.
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type key int

const dbKey key = 0

func WithExecutor(ctx context.Context, db Executor) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

func ExecutorFromContext(ctx context.Context) Executor {
	db, ok := ctx.Value(dbKey).(Executor)
	if !ok {
		panic("No database in the context")
	}
	return db
}

func ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	exec := ExecutorFromContext(ctx)
	return exec.ExecContext(ctx, query, args...)
}
func QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	exec := ExecutorFromContext(ctx)
	return exec.QueryContext(ctx, query, args...)
}
func QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	exec := ExecutorFromContext(ctx)
	return exec.QueryRowContext(ctx, query, args...)
}

// Atomic invokes the passed function in the context of a managed SQL
// transaction.  Any errors returned from the user-supplied function are
// returned from this function.
//
// Retries are automatically performed in case of serialization failures or deadlocks.
func Atomic(ctx context.Context, fn func(ctx context.Context) error) error {
	var err error
	for try := 0; try < 5; try++ {
		err = doTransaction(ctx, fn, false)
		if err == nil || !shouldRetryTransaction(err) {
			return err
		}
	}
	return err
}

// AtomicReadOnly invokes the passed function in the context of a managed SQL
// read only transaction. Any errors returned from the user-supplied function
// are returned from this function.
//
// Retries are automatically performed in case of serialization failures or deadlocks.
func AtomicReadOnly(ctx context.Context, fn func(ctx context.Context) error) error {
	var err error
	for try := 0; try < 5; try++ {
		err = doTransaction(ctx, fn, true)
		if err == nil || !shouldRetryTransaction(err) {
			return err
		}
	}
	return err
}

func shouldRetryTransaction(err error) bool {
	err2 := errors.Cause(err)
	if err2, ok := err2.(*pq.Error); ok {
		n := err2.Code.Name()
		if n == "serialization_failure" || n == "deadlock_detected" {
			return true
		}
	}
	return false
}

// rollbackOnPanic rolls the passed transaction back if the code in the calling
// function panics. This is needed in order to not leak transactions in case
// of panic.
func rollbackOnPanic(tx *sql.Tx) {
	if err := recover(); err != nil {
		_ = tx.Rollback()
		panic(err)
	}
}

// Transaction invokes the passed function in the context of a managed SQL
// transaction.  Any errors returned from
// the user-supplied function are returned from this function.
func doTransaction(ctx context.Context, fn func(ctx context.Context) error, readOnly bool) error {
	db, ok := ExecutorFromContext(ctx).(*sql.DB)
	if !ok {
		panic("database does not support transactions")
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  readOnly,
	})
	if err != nil {
		// TODO rollback??
		return err
	}

	ctx2 := WithExecutor(ctx, tx)

	// Since the user-provided function might panic, ensure the transaction
	// releases all resources.
	defer rollbackOnPanic(tx)

	err = fn(ctx2)
	if err != nil {
		// Error ignored here, maybe we should do something with it?
		// Not sure what though.
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		// TODO rollback?
		return err
	}
	return nil
}
