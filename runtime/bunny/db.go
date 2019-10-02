package bunny

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/lib/pq"
	"github.com/sqlbunny/errors"
)

// DB can perform SQL queries.
type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type key int

const dbKey key = 0

func ContextWithDB(ctx context.Context, db DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

func DBFromContext(ctx context.Context) DB {
	db, ok := ctx.Value(dbKey).(DB)
	if !ok {
		panic("No database in the context")
	}
	return db
}

func Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	db := DBFromContext(ctx)
	begin := time.Now()
	res, err := db.ExecContext(ctx, query, args...)
	if logger != nil {
		logger.LogQuery(ctx, QueryLogInfo{
			Query:    query,
			Duration: time.Since(begin),
			Err:      err,
			Args:     args,
		})
	}
	return res, err
}

func Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	db := DBFromContext(ctx)
	begin := time.Now()
	res, err := db.QueryContext(ctx, query, args...)
	if logger != nil {
		logger.LogQuery(ctx, QueryLogInfo{
			Query:    query,
			Duration: time.Since(begin),
			Err:      err,
			Args:     args,
		})
	}
	return res, err
}

func QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	db := DBFromContext(ctx)
	begin := time.Now()
	res := db.QueryRowContext(ctx, query, args...)
	if logger != nil {
		logger.LogQuery(ctx, QueryLogInfo{
			Query:    query,
			Duration: time.Since(begin),
			Err:      nil, // TODO how to get the error without causing quantum decoherence in res?
			Args:     args,
		})
	}
	return res
}

// Atomic invokes the passed function in the context of a managed SQL
// transaction.  Any errors returned from the user-supplied function are
// returned from this function.
//
// Retries are automatically performed in case of serialization failures or deadlocks.
func Atomic(ctx context.Context, fn func(ctx context.Context) error) error {
	return doAtomic(ctx, fn, false)
}

// AtomicReadOnly invokes the passed function in the context of a managed SQL
// read only transaction. Any errors returned from the user-supplied function
// are returned from this function.
//
// Retries are automatically performed in case of serialization failures or deadlocks.
func AtomicReadOnly(ctx context.Context, fn func(ctx context.Context) error) error {
	return doAtomic(ctx, fn, true)
}

func doAtomic(ctx context.Context, fn func(ctx context.Context) error, readOnly bool) error {
	var err error
	for try := uint(0); try < 12; try++ {
		err = doTransaction(ctx, fn, readOnly)
		if err == nil {
			return nil
		}

		if !shouldRetryTransaction(err) {
			return err
		}

		time.Sleep(time.Millisecond * time.Duration(rand.Int31n(1<<try)))
	}
	return err
}

func shouldRetryTransaction(err error) bool {
	err2 := errors.Unwrap(err)
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
func rollbackOnPanic(ctx context.Context, tx *txNode, begin time.Time) {
	if err := recover(); err != nil {
		if logger != nil {
			logger.LogRollback(ctx, RollbackLogInfo{
				Duration: time.Since(begin),
				Err:      fmt.Errorf("panic %v", err),
			})
		}

		err2 := tx.Rollback()
		if err2 != nil {
			panic(err2)
		}
		panic(err)
	}
}

type txNode struct {
	dbTx   *sql.Tx
	parent *txNode
	child  *txNode
	depth  int
}

func (t *txNode) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if t.child != nil {
		panic("Transaction has a subtransaction active, can't run statements in it.")
	}
	return t.dbTx.ExecContext(ctx, query, args...)
}
func (t *txNode) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if t.child != nil {
		panic("Transaction has a subtransaction active, can't run statements in it.")
	}
	return t.dbTx.QueryContext(ctx, query, args...)
}
func (t *txNode) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if t.child != nil {
		panic("Transaction has a subtransaction active, can't run statements in it.")
	}
	return t.dbTx.QueryRowContext(ctx, query, args...)
}

func (t *txNode) Commit() error {
	if t.parent == nil {
		return t.dbTx.Commit()
	}
	_, err := t.dbTx.Exec(fmt.Sprintf("RELEASE SAVEPOINT savepoint_%d", t.depth))
	t.parent.child = nil
	return err
}

func (t *txNode) Rollback() error {
	if t.parent == nil {
		return t.dbTx.Rollback()
	}
	_, err := t.dbTx.Exec(fmt.Sprintf("ROLLBACK TO SAVEPOINT savepoint_%d", t.depth))
	t.parent.child = nil
	return err
}

type beginTxer interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// Transaction invokes the passed function in the context of a managed SQL
// transaction.  Any errors returned from
// the user-supplied function are returned from this function.
func doTransaction(ctx context.Context, fn func(ctx context.Context) error, readOnly bool) error {
	if logger != nil {
		ctx = logger.LogBegin(ctx, BeginLogInfo{
			ReadOnly: readOnly,
		})
	}
	begin := time.Now()

	var node *txNode
	switch db := DBFromContext(ctx).(type) {
	case beginTxer:
		tx, err := db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  readOnly,
		})
		if err != nil {
			retErr := errors.Errorf("BeginTx failed: %w", err)
			if logger != nil {
				logger.LogRollback(ctx, RollbackLogInfo{
					Duration: time.Since(begin),
					Err:      retErr,
				})
			}
			return retErr
		}
		node = &txNode{
			dbTx:  tx,
			depth: 0,
		}
	case *txNode:
		node = &txNode{
			dbTx:   db.dbTx,
			parent: db,
			depth:  db.depth + 1,
		}
		_, err := db.dbTx.Exec(fmt.Sprintf("SAVEPOINT savepoint_%d", node.depth))
		if err != nil {
			return err
		}
		db.child = node
	default:
		panic("database does not support transactions")
	}

	ctx2 := ContextWithDB(ctx, node)

	// Since the user-provided function might panic, ensure the transaction
	// releases all resources.
	defer rollbackOnPanic(ctx, node, begin)

	err := fn(ctx2)
	if err != nil {
		retErr := errors.Errorf("tx function returned error: %w", err)
		if logger != nil {
			logger.LogRollback(ctx, RollbackLogInfo{
				Duration: time.Since(begin),
				Err:      retErr,
			})
		}
		err2 := node.Rollback()
		if err2 != nil {
			panic(err2)
		}
		return retErr
	}

	err = node.Commit()
	if err != nil {
		retErr := errors.Errorf("commit: %w", err)
		if logger != nil {
			logger.LogRollback(ctx, RollbackLogInfo{
				Duration: time.Since(begin),
				Err:      retErr,
			})
		}
		_ = node.Rollback()
		return retErr
	}
	if logger != nil {
		logger.LogCommit(ctx, CommitLogInfo{
			Duration: time.Since(begin),
		})
	}

	return nil
}

func IsAtomic(ctx context.Context) bool {
	_, ok := DBFromContext(ctx).(*txNode)
	return ok
}

func AssertAtomic(ctx context.Context) {
	if !IsAtomic(ctx) {
		panic("AssertAtomic failed")
	}
}

func AssertNotAtomic(ctx context.Context) {
	if IsAtomic(ctx) {
		panic("AssertNotAtomic failed")
	}
}
