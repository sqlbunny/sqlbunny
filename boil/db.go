package boil

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	begin := time.Now()
	res, err := exec.ExecContext(ctx, query, args...)
	if queryLogger != nil {
		queryLogger.LogQuery(ctx, query, time.Since(begin), args...)
	}
	return res, err
}

func QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	exec := ExecutorFromContext(ctx)
	begin := time.Now()
	res, err := exec.QueryContext(ctx, query, args...)
	if queryLogger != nil {
		queryLogger.LogQuery(ctx, query, time.Since(begin), args...)
	}
	return res, err
}

func QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	exec := ExecutorFromContext(ctx)
	begin := time.Now()
	res := exec.QueryRowContext(ctx, query, args...)
	if queryLogger != nil {
		queryLogger.LogQuery(ctx, query, time.Since(begin), args...)
	}
	return res
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
func rollbackOnPanic(ctx context.Context, tx *sql.Tx, begin time.Time) {
	if err := recover(); err != nil {
		if queryLogger != nil {
			queryLogger.LogRollback(ctx, time.Since(begin), fmt.Errorf("panic %v", err))
		}

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

	if queryLogger != nil {
		ctx = queryLogger.LogBegin(ctx, readOnly)
	}
	begin := time.Now()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  readOnly,
	})
	if err != nil {
		if queryLogger != nil {
			queryLogger.LogRollback(ctx, time.Since(begin), errors.Wrap(err, "BeginTx:"))
		}
		_ = tx.Rollback()
		return err
	}

	ctx2 := WithExecutor(ctx, tx)

	// Since the user-provided function might panic, ensure the transaction
	// releases all resources.
	defer rollbackOnPanic(ctx, tx, begin)

	err = fn(ctx2)
	if err != nil {
		// Error ignored here, maybe we should do something with it?
		// Not sure what though.
		if queryLogger != nil {
			queryLogger.LogRollback(ctx, time.Since(begin), errors.Wrap(err, "tx function returned error:"))
		}
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		if queryLogger != nil {
			queryLogger.LogRollback(ctx, time.Since(begin), errors.Wrap(err, "Commit:"))
		}
		_ = tx.Rollback()
		return err
	}
	if queryLogger != nil {
		queryLogger.LogCommit(ctx, time.Since(begin))
	}

	return nil
}
