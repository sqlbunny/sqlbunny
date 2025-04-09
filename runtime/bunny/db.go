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
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type contextDBKeyType struct{}

var ContextDBKey = contextDBKeyType{}

func ContextWithDB(ctx context.Context, db DB) context.Context {
	return context.WithValue(ctx, ContextDBKey, db)
}

func DBFromContext(ctx context.Context) DB {
	db, ok := ctx.Value(ContextDBKey).(DB)
	if !ok {
		panic("No database in the context")
	}
	return db
}

func Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	db := DBFromContext(ctx)
	begin := time.Now()
	res, err := db.ExecContext(ctx, query, args...)
	err = errors.WithStack(err)
	logger.LogQuery(ctx, QueryLogInfo{
		Query:    query,
		Duration: time.Since(begin),
		Err:      err,
		Args:     args,
	})
	return res, err
}

func Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	db := DBFromContext(ctx)
	begin := time.Now()
	res, err := db.QueryContext(ctx, query, args...)
	err = errors.WithStack(err)
	logger.LogQuery(ctx, QueryLogInfo{
		Query:    query,
		Duration: time.Since(begin),
		Err:      err,
		Args:     args,
	})
	return res, err
}

func QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	db := DBFromContext(ctx)
	begin := time.Now()
	res := db.QueryRowContext(ctx, query, args...)
	logger.LogQuery(ctx, QueryLogInfo{
		Query:    query,
		Duration: time.Since(begin),
		Err:      nil, // TODO how to get the error without causing quantum decoherence in res?
		Args:     args,
	})
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
	var pqerr *pq.Error
	if errors.As(err, &pqerr) {
		switch pqerr.Code {
		case "40001": // serialization_failures
			return true
		case "40P01": // deadlock_detected
			return true
		}
	}
	return false
}

type txNode struct {
	dbTx     *sql.Tx
	parent   *txNode
	child    *txNode
	depth    int
	onCommit []func(context.Context) error
}

func (t *txNode) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if t.child != nil {
		panic("Transaction has a subtransaction active, can't run statements in it.")
	}
	return t.dbTx.ExecContext(ctx, query, args...)
}
func (t *txNode) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if t.child != nil {
		panic("Transaction has a subtransaction active, can't run statements in it.")
	}
	return t.dbTx.QueryContext(ctx, query, args...)
}
func (t *txNode) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
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

func (t *txNode) runOnCommit(ctx context.Context) error {
	if t.parent != nil {
		// If we're in a subtransaction, we don't want to execute the onCommits yet.
		// We append them to the parent transaction, so they'll run when the topmost transaction commits.
		t.parent.onCommit = append(t.parent.onCommit, t.onCommit...)
		return nil
	}

	// We are the top most transaction.
	// Just run the onCommit hooks.
	for _, fn := range t.onCommit {
		err := fn(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

type beginTxer interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// Transaction invokes the passed function in the context of a managed SQL
// transaction.  Any errors returned from
// the user-supplied function are returned from this function.
func doTransaction(ctx context.Context, fn func(ctx context.Context) error, readOnly bool) (err error) {
	ctx = logger.LogBegin(ctx, BeginLogInfo{
		ReadOnly: readOnly,
	})
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
			logger.LogRollback(ctx, RollbackLogInfo{
				Duration: time.Since(begin),
				Err:      retErr,
			})
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

	committed := false
	// Since the user-provided function might panic, ensure the transaction
	// releases all resources.
	defer func() {
		if committed {
			return
		}

		if p := recover(); p != nil {
			logger.LogRollback(ctx, RollbackLogInfo{
				Duration: time.Since(begin),
				Err:      fmt.Errorf("panic %v", p),
			})
			_ = node.Rollback() // just ignore errors here.
			panic(p)            // re-throw panic.
		} else if err != nil {
			logger.LogRollback(ctx, RollbackLogInfo{
				Duration: time.Since(begin),
				Err:      err,
			})
			err2 := node.Rollback()
			if err2 != nil {
				err = errors.Errorf("tx failed: %w, and rollback failed: %w", err, err2)
			}
		}
	}()

	err = fn(ctx2)
	if err != nil {
		return errors.Errorf("tx function returned error: %w", err)
	}

	err = node.Commit()
	if err != nil {
		return errors.Errorf("commit: %w", err)
	}
	committed = true
	logger.LogCommit(ctx, CommitLogInfo{
		Duration: time.Since(begin),
	})

	err = node.runOnCommit(ctx)
	if err != nil {
		return errors.Errorf("OnCommit returned error: %w", err)
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

func OnCommit(ctx context.Context, fn func(context.Context) error) {
	db := DBFromContext(ctx)
	tx, ok := db.(*txNode)
	if !ok {
		panic("OnCommit called while not in atomic")
	}

	tx.onCommit = append(tx.onCommit, fn)
}
