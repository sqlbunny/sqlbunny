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

// AtomicOptions configures transaction behavior for AtomicWithOptions.
type AtomicOptions struct {
	ReadOnly  bool
	Isolation sql.IsolationLevel // 0 means use default (Serializable)
}

// Atomic invokes the passed function in the context of a managed SQL
// transaction.  Any errors returned from the user-supplied function are
// returned from this function.
//
// Retries are automatically performed in case of serialization failures or deadlocks.
func Atomic(ctx context.Context, fn func(ctx context.Context) error) error {
	return doAtomic(ctx, fn, AtomicOptions{})
}

// AtomicReadOnly invokes the passed function in the context of a managed SQL
// read only transaction. Any errors returned from the user-supplied function
// are returned from this function.
//
// Retries are automatically performed in case of serialization failures or deadlocks.
func AtomicReadOnly(ctx context.Context, fn func(ctx context.Context) error) error {
	return doAtomic(ctx, fn, AtomicOptions{ReadOnly: true})
}

// AtomicWithOptions invokes the passed function in the context of a managed SQL
// transaction with the given options.
//
// Retries are automatically performed in case of serialization failures or deadlocks.
func AtomicWithOptions(ctx context.Context, opts AtomicOptions, fn func(ctx context.Context) error) error {
	return doAtomic(ctx, fn, opts)
}

func doAtomic(ctx context.Context, fn func(ctx context.Context) error, opts AtomicOptions) error {
	// If we're already inside a transaction, don't retry here.
	// Serialization failures poison the entire transaction, so retrying
	// a savepoint is pointless. Let the error propagate to the outermost
	// Atomic which can retry the whole transaction.
	if _, ok := DBFromContext(ctx).(*Tx); ok {
		return doTransaction(ctx, fn, opts)
	}

	var err error
	for try := uint(0); try < 12; try++ {
		err = doTransaction(ctx, fn, opts)
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

// Tx represents one level of a managed transaction: either the root
// transaction (depth=0, an actual BEGIN/COMMIT) or a savepoint inside a
// parent transaction (depth>0, a SAVEPOINT/RELEASE pair).
//
// Tx exposes per-level hooks (BeforeCommit, AfterCommit) that fire when
// this specific level commits — i.e. when this savepoint is released, or
// when the root transaction's SQL COMMIT lands. Hooks registered on a
// subtx whose savepoint gets rolled back are silently discarded along
// with the subtx; they only fire if the subtx successfully commits.
//
// Tx also exposes a per-level mutable state map (GetState, SetState) for
// components that want to attach metadata that should die with the subtx
// if it rolls back. State on a parent is visible only to that parent
// (not to descendants); use CurrentTx(ctx).Parent() to reach it.
type Tx struct {
	dbTx         *sql.Tx
	parent       *Tx
	child        *Tx
	depth        int
	beforeCommit []func(ctx context.Context, tx *Tx) error
	afterCommit  []func(ctx context.Context, tx *Tx) error
	state        map[any]any
}

// Parent returns this Tx's parent, or nil if this is the root transaction.
func (t *Tx) Parent() *Tx { return t.parent }

// BeforeCommit registers fn to run just before this subtx commits (i.e.
// before SAVEPOINT RELEASE for a subtx, before SQL COMMIT for the root).
// fn runs *inside* the transaction, so its DB writes are part of the same
// atomic unit and are rolled back together if either fn or the commit
// itself fails.
//
// Hooks fire in registration order and may register additional hooks; new
// hooks added during iteration are picked up in the same loop.
//
// To match the semantics of the legacy package-level BeforeCommit (run
// only when the root transaction commits, never on savepoint release),
// register a self-bubbling closure — see the BeforeCommit package
// function for a ready-made implementation.
func (t *Tx) BeforeCommit(fn func(ctx context.Context, tx *Tx) error) {
	t.beforeCommit = append(t.beforeCommit, fn)
}

// AfterCommit registers fn to run just after this subtx commits (i.e.
// after SAVEPOINT RELEASE for a subtx, after SQL COMMIT for the root).
// For a subtx, AfterCommit fires inside the outer transaction so DB
// writes still work; for the root tx it fires after the actual commit,
// outside any transaction.
//
// Hooks added on a subtx that rolls back are discarded.
func (t *Tx) AfterCommit(fn func(ctx context.Context, tx *Tx) error) {
	t.afterCommit = append(t.afterCommit, fn)
}

// GetState returns the value stored on this Tx under key, or nil if no
// value was set. State is per-Tx (not shared with parent / children); use
// Parent() to reach into the enclosing tx's state.
func (t *Tx) GetState(key any) any {
	return t.state[key]
}

// SetState stores val on this Tx under key. Passing nil deletes the entry.
// State lives only as long as this Tx — when a subtx rolls back, its state
// is discarded along with its hooks.
func (t *Tx) SetState(key, val any) {
	if val == nil {
		delete(t.state, key)
		return
	}
	if t.state == nil {
		t.state = make(map[any]any)
	}
	t.state[key] = val
}

// CurrentTx returns the innermost Tx for ctx, or nil if ctx is not inside
// a transaction. Use this to register per-level hooks or attach per-level
// state from libraries that build on top of bunny.
func CurrentTx(ctx context.Context) *Tx {
	tx, _ := ctx.Value(ContextDBKey).(*Tx)
	return tx
}

func (t *Tx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if t.child != nil {
		panic("Transaction has a subtransaction active, can't run statements in it.")
	}
	return t.dbTx.ExecContext(ctx, query, args...)
}
func (t *Tx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if t.child != nil {
		panic("Transaction has a subtransaction active, can't run statements in it.")
	}
	return t.dbTx.QueryContext(ctx, query, args...)
}
func (t *Tx) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if t.child != nil {
		panic("Transaction has a subtransaction active, can't run statements in it.")
	}
	return t.dbTx.QueryRowContext(ctx, query, args...)
}

func (t *Tx) commit() error {
	if t.parent == nil {
		return t.dbTx.Commit()
	}
	_, err := t.dbTx.Exec(fmt.Sprintf("RELEASE SAVEPOINT savepoint_%d", t.depth))
	t.parent.child = nil
	return err
}

func (t *Tx) rollback() error {
	if t.parent == nil {
		return t.dbTx.Rollback()
	}
	_, err := t.dbTx.Exec(fmt.Sprintf("ROLLBACK TO SAVEPOINT savepoint_%d", t.depth))
	t.parent.child = nil
	return err
}

// runHooks runs each hook in `hooks` in order, with this Tx as the
// argument, and using an index loop so hooks added during iteration are
// also picked up. Resets the hook slice when done.
func (t *Tx) runHooks(ctx context.Context, which *[]func(ctx context.Context, tx *Tx) error) error {
	for i := 0; i < len(*which); i++ {
		if err := (*which)[i](ctx, t); err != nil {
			return err
		}
	}
	*which = nil
	return nil
}

type beginTxer interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// Transaction invokes the passed function in the context of a managed SQL
// transaction.  Any errors returned from
// the user-supplied function are returned from this function.
func doTransaction(ctx context.Context, fn func(ctx context.Context) error, opts AtomicOptions) (err error) {
	ctx = logger.LogBegin(ctx, BeginLogInfo{
		ReadOnly: opts.ReadOnly,
	})
	begin := time.Now()

	isolation := opts.Isolation
	if isolation == 0 {
		isolation = sql.LevelSerializable
	}

	var node *Tx
	switch db := DBFromContext(ctx).(type) {
	case beginTxer:
		tx, err := db.BeginTx(ctx, &sql.TxOptions{
			Isolation: isolation,
			ReadOnly:  opts.ReadOnly,
		})
		if err != nil {
			retErr := errors.Errorf("BeginTx failed: %w", err)
			logger.LogRollback(ctx, RollbackLogInfo{
				Duration: time.Since(begin),
				Err:      retErr,
			})
			return retErr
		}
		node = &Tx{
			dbTx:  tx,
			depth: 0,
		}
	case *Tx:
		node = &Tx{
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
			_ = node.rollback() // just ignore errors here.
			panic(p)            // re-throw panic.
		} else if err != nil {
			logger.LogRollback(ctx, RollbackLogInfo{
				Duration: time.Since(begin),
				Err:      err,
			})
			err2 := node.rollback()
			if err2 != nil {
				err = errors.Errorf("tx failed: %w, and rollback failed: %w", err, err2)
			}
		}
	}()

	err = fn(ctx2)
	if err != nil {
		return errors.Errorf("tx function returned error: %w", err)
	}

	err = node.runHooks(ctx2, &node.beforeCommit)
	if err != nil {
		return errors.Errorf("BeforeCommit returned error: %w", err)
	}

	err = node.commit()
	if err != nil {
		return errors.Errorf("commit: %w", err)
	}
	committed = true
	logger.LogCommit(ctx, CommitLogInfo{
		Duration: time.Since(begin),
	})

	// AfterCommit fires after the commit. For a subtx (savepoint release)
	// the parent tx is still open and ctx already carries it (that's how
	// we picked up the parent in doTransaction's switch). For the root,
	// ctx carries no tx — correct, since AfterCommit fires after COMMIT.
	err = node.runHooks(ctx, &node.afterCommit)
	if err != nil {
		return errors.Errorf("AfterCommit returned error: %w", err)
	}

	return nil
}

func IsAtomic(ctx context.Context) bool {
	return CurrentTx(ctx) != nil
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
