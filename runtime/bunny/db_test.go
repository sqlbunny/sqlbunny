package bunny

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/sqlbunny/errors"
	"gopkg.in/DATA-DOG/go-sqlmock.v2"
)

// recorder collects events from both mock SQL expectations and user code.
type recorder struct {
	events []string
}

func (r *recorder) record(event string) {
	r.events = append(r.events, event)
}

func (r *recorder) check(t *testing.T, expected []string) {
	t.Helper()
	got := strings.Join(r.events, "\n")
	want := strings.Join(expected, "\n")
	if got != want {
		t.Errorf("events mismatch\ngot:\n%s\n\nwant:\n%s", got, want)
	}
}

var errTest = fmt.Errorf("test error")

var errSerialization = &pq.Error{Code: "40001"}

var errDeadlock = &pq.Error{Code: "40P01"}

func setupTest(t *testing.T) (context.Context, sqlmock.Sqlmock, *recorder) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	ctx := ContextWithDB(context.Background(), db)
	rec := &recorder{}
	return ctx, mock, rec
}

func TestAtomic_Success(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectCommit()

	err := Atomic(ctx, func(ctx context.Context) error {
		rec.record("fn")
		return nil
	})

	rec.record(fmt.Sprintf("err=%v", err))
	rec.check(t, []string{
		"fn",
		"err=<nil>",
	})
}

func TestAtomic_FnError(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectRollback()

	err := Atomic(ctx, func(ctx context.Context) error {
		rec.record("fn")
		return errTest
	})

	rec.record(fmt.Sprintf("err is errTest=%v", errors.Is(err, errTest)))
	rec.check(t, []string{
		"fn",
		"err is errTest=true",
	})
}

func TestAtomic_SerializationRetry(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	// First attempt: fn returns serialization error -> retry
	mock.ExpectBegin()
	mock.ExpectRollback()
	// Second attempt: succeeds
	mock.ExpectBegin()
	mock.ExpectCommit()

	attempt := 0
	err := Atomic(ctx, func(ctx context.Context) error {
		attempt++
		rec.record(fmt.Sprintf("fn attempt=%d", attempt))
		if attempt == 1 {
			return errSerialization
		}
		return nil
	})

	rec.record(fmt.Sprintf("err=%v", err))
	rec.check(t, []string{
		"fn attempt=1",
		"fn attempt=2",
		"err=<nil>",
	})
}

func TestAtomic_DeadlockRetry(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectRollback()
	mock.ExpectBegin()
	mock.ExpectCommit()

	attempt := 0
	err := Atomic(ctx, func(ctx context.Context) error {
		attempt++
		rec.record(fmt.Sprintf("fn attempt=%d", attempt))
		if attempt == 1 {
			return errDeadlock
		}
		return nil
	})

	rec.record(fmt.Sprintf("err=%v", err))
	rec.check(t, []string{
		"fn attempt=1",
		"fn attempt=2",
		"err=<nil>",
	})
}

func TestAtomic_SerializationExhausted(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	for i := 0; i < 12; i++ {
		mock.ExpectBegin()
		mock.ExpectRollback()
	}

	attempt := 0
	err := Atomic(ctx, func(ctx context.Context) error {
		attempt++
		rec.record(fmt.Sprintf("fn attempt=%d", attempt))
		return errSerialization
	})

	var pqErr *pq.Error
	isPQ := errors.As(err, &pqErr)
	rec.record(fmt.Sprintf("isPQ=%v code=%s", isPQ, pqErr.Code))
	rec.check(t, []string{
		"fn attempt=1",
		"fn attempt=2",
		"fn attempt=3",
		"fn attempt=4",
		"fn attempt=5",
		"fn attempt=6",
		"fn attempt=7",
		"fn attempt=8",
		"fn attempt=9",
		"fn attempt=10",
		"fn attempt=11",
		"fn attempt=12",
		"isPQ=true code=40001",
	})
}

func TestAtomic_NonRetryableError(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectRollback()

	err := Atomic(ctx, func(ctx context.Context) error {
		rec.record("fn")
		return errTest
	})

	rec.record(fmt.Sprintf("err is errTest=%v", errors.Is(err, errTest)))
	rec.check(t, []string{
		"fn",
		"err is errTest=true",
	})
}

func TestAtomic_ReadOnly(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectCommit()

	err := AtomicReadOnly(ctx, func(ctx context.Context) error {
		rec.record("fn")
		return nil
	})

	rec.record(fmt.Sprintf("err=%v", err))
	rec.check(t, []string{
		"fn",
		"err=<nil>",
	})
}

func TestAtomic_WithOptions_Isolation(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectCommit()

	err := AtomicWithOptions(ctx, AtomicOptions{Isolation: sql.LevelRepeatableRead}, func(ctx context.Context) error {
		rec.record("fn")
		return nil
	})

	rec.record(fmt.Sprintf("err=%v", err))
	rec.check(t, []string{
		"fn",
		"err=<nil>",
	})
}

func TestAtomic_OnCommit(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectCommit()

	err := Atomic(ctx, func(ctx context.Context) error {
		rec.record("fn")
		OnCommit(ctx, func(ctx context.Context) error {
			rec.record("oncommit1")
			return nil
		})
		OnCommit(ctx, func(ctx context.Context) error {
			rec.record("oncommit2")
			return nil
		})
		return nil
	})

	rec.record(fmt.Sprintf("err=%v", err))
	rec.check(t, []string{
		"fn",
		"oncommit1",
		"oncommit2",
		"err=<nil>",
	})
}

func TestAtomic_OnCommit_NotCalledOnError(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectRollback()

	err := Atomic(ctx, func(ctx context.Context) error {
		rec.record("fn")
		OnCommit(ctx, func(ctx context.Context) error {
			rec.record("oncommit")
			return nil
		})
		return errTest
	})

	rec.record(fmt.Sprintf("err is errTest=%v", errors.Is(err, errTest)))
	rec.check(t, []string{
		"fn",
		"err is errTest=true",
	})
}

func TestAtomic_OnCommit_Error(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectCommit()

	err := Atomic(ctx, func(ctx context.Context) error {
		rec.record("fn")
		OnCommit(ctx, func(ctx context.Context) error {
			rec.record("oncommit1")
			return errTest
		})
		OnCommit(ctx, func(ctx context.Context) error {
			rec.record("oncommit2")
			return nil
		})
		return nil
	})

	rec.record(fmt.Sprintf("err is errTest=%v", errors.Is(err, errTest)))
	rec.check(t, []string{
		"fn",
		"oncommit1",
		"err is errTest=true",
	})
}

// Nested transactions (savepoints)

func TestAtomic_Nested_Success(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectExec("SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("RELEASE SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := Atomic(ctx, func(ctx context.Context) error {
		rec.record("outer")
		err := Atomic(ctx, func(ctx context.Context) error {
			rec.record("inner")
			return nil
		})
		rec.record(fmt.Sprintf("inner err=%v", err))
		return nil
	})

	rec.record(fmt.Sprintf("err=%v", err))
	rec.check(t, []string{
		"outer",
		"inner",
		"inner err=<nil>",
		"err=<nil>",
	})
}

func TestAtomic_Nested_InnerError(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectExec("SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("ROLLBACK TO SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := Atomic(ctx, func(ctx context.Context) error {
		rec.record("outer")
		innerErr := Atomic(ctx, func(ctx context.Context) error {
			rec.record("inner")
			return errTest
		})
		rec.record(fmt.Sprintf("inner err is errTest=%v", errors.Is(innerErr, errTest)))
		return nil
	})

	rec.record(fmt.Sprintf("err=%v", err))
	rec.check(t, []string{
		"outer",
		"inner",
		"inner err is errTest=true",
		"err=<nil>",
	})
}

func TestAtomic_Nested_InnerError_Propagated(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectExec("SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("ROLLBACK TO SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := Atomic(ctx, func(ctx context.Context) error {
		rec.record("outer")
		innerErr := Atomic(ctx, func(ctx context.Context) error {
			rec.record("inner")
			return errTest
		})
		rec.record(fmt.Sprintf("inner err is errTest=%v", errors.Is(innerErr, errTest)))
		return innerErr
	})

	rec.record(fmt.Sprintf("err is errTest=%v", errors.Is(err, errTest)))
	rec.check(t, []string{
		"outer",
		"inner",
		"inner err is errTest=true",
		"err is errTest=true",
	})
}

func TestAtomic_Nested_SerializationNotRetried(t *testing.T) {
	// Serialization errors in nested Atomic should NOT be retried,
	// because the entire parent transaction is poisoned.
	ctx, mock, rec := setupTest(t)

	// Outer attempt 1
	mock.ExpectBegin()
	mock.ExpectExec("SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("ROLLBACK TO SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()
	// Outer attempt 2 (retry)
	mock.ExpectBegin()
	mock.ExpectExec("SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("RELEASE SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	outerAttempt := 0
	err := Atomic(ctx, func(ctx context.Context) error {
		outerAttempt++
		rec.record(fmt.Sprintf("outer attempt=%d", outerAttempt))
		innerErr := Atomic(ctx, func(ctx context.Context) error {
			rec.record(fmt.Sprintf("inner attempt=%d", outerAttempt))
			if outerAttempt == 1 {
				return errSerialization
			}
			return nil
		})
		if innerErr != nil {
			return innerErr
		}
		return nil
	})

	rec.record(fmt.Sprintf("err=%v", err))
	rec.check(t, []string{
		"outer attempt=1",
		"inner attempt=1",
		"outer attempt=2",
		"inner attempt=2",
		"err=<nil>",
	})
}

func TestAtomic_Nested_DeadlockNotRetried(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	// Outer attempt 1
	mock.ExpectBegin()
	mock.ExpectExec("SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("ROLLBACK TO SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()
	// Outer attempt 2 (retry)
	mock.ExpectBegin()
	mock.ExpectExec("SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("RELEASE SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	outerAttempt := 0
	err := Atomic(ctx, func(ctx context.Context) error {
		outerAttempt++
		rec.record(fmt.Sprintf("outer attempt=%d", outerAttempt))
		innerErr := Atomic(ctx, func(ctx context.Context) error {
			rec.record(fmt.Sprintf("inner attempt=%d", outerAttempt))
			if outerAttempt == 1 {
				return errDeadlock
			}
			return nil
		})
		if innerErr != nil {
			return innerErr
		}
		return nil
	})

	rec.record(fmt.Sprintf("err=%v", err))
	rec.check(t, []string{
		"outer attempt=1",
		"inner attempt=1",
		"outer attempt=2",
		"inner attempt=2",
		"err=<nil>",
	})
}

func TestAtomic_Nested_OnCommit_DefersToOuter(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectExec("SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("RELEASE SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := Atomic(ctx, func(ctx context.Context) error {
		rec.record("outer")
		OnCommit(ctx, func(ctx context.Context) error {
			rec.record("outer oncommit")
			return nil
		})
		err := Atomic(ctx, func(ctx context.Context) error {
			rec.record("inner")
			OnCommit(ctx, func(ctx context.Context) error {
				rec.record("inner oncommit")
				return nil
			})
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})

	rec.record(fmt.Sprintf("err=%v", err))
	rec.check(t, []string{
		"outer",
		"inner",
		"outer oncommit",
		"inner oncommit",
		"err=<nil>",
	})
}

func TestAtomic_Nested_OnCommit_NotCalledOnInnerError(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectExec("SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("ROLLBACK TO SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := Atomic(ctx, func(ctx context.Context) error {
		rec.record("outer")
		OnCommit(ctx, func(ctx context.Context) error {
			rec.record("outer oncommit")
			return nil
		})
		_ = Atomic(ctx, func(ctx context.Context) error {
			rec.record("inner")
			OnCommit(ctx, func(ctx context.Context) error {
				rec.record("inner oncommit")
				return nil
			})
			return errTest
		})
		// outer still succeeds
		return nil
	})

	rec.record(fmt.Sprintf("err=%v", err))
	rec.check(t, []string{
		"outer",
		"inner",
		"outer oncommit",
		"err=<nil>",
	})
}

func TestAtomic_Nested_TwoLevels(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectExec("SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("SAVEPOINT savepoint_2").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("RELEASE SAVEPOINT savepoint_2").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("RELEASE SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := Atomic(ctx, func(ctx context.Context) error {
		rec.record("level0")
		return Atomic(ctx, func(ctx context.Context) error {
			rec.record("level1")
			return Atomic(ctx, func(ctx context.Context) error {
				rec.record("level2")
				return nil
			})
		})
	})

	rec.record(fmt.Sprintf("err=%v", err))
	rec.check(t, []string{
		"level0",
		"level1",
		"level2",
		"err=<nil>",
	})
}

func TestAtomic_Panic(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectRollback()

	defer func() {
		p := recover()
		rec.record(fmt.Sprintf("recovered=%v", p))
		rec.check(t, []string{
			"fn",
			"recovered=boom",
		})
	}()

	_ = Atomic(ctx, func(ctx context.Context) error {
		rec.record("fn")
		panic("boom")
	})
}

func TestAtomic_Nested_Panic(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectExec("SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("ROLLBACK TO SAVEPOINT savepoint_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	defer func() {
		p := recover()
		rec.record(fmt.Sprintf("recovered=%v", p))
		rec.check(t, []string{
			"outer",
			"inner",
			"recovered=boom",
		})
	}()

	_ = Atomic(ctx, func(ctx context.Context) error {
		rec.record("outer")
		return Atomic(ctx, func(ctx context.Context) error {
			rec.record("inner")
			panic("boom")
		})
	})
}

func TestAtomic_CommitError(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(errTest)
	mock.ExpectRollback()

	err := Atomic(ctx, func(ctx context.Context) error {
		rec.record("fn")
		return nil
	})

	rec.record(fmt.Sprintf("err contains test error=%v", errors.Is(err, errTest)))
	rec.check(t, []string{
		"fn",
		"err contains test error=true",
	})
}

func TestAtomic_BeginError(t *testing.T) {
	ctx, mock, rec := setupTest(t)

	mock.ExpectBegin().WillReturnError(errTest)

	err := Atomic(ctx, func(ctx context.Context) error {
		rec.record("fn")
		return nil
	})

	rec.record(fmt.Sprintf("fn called=%v", len(rec.events) > 0 && rec.events[0] == "fn"))
	rec.record(fmt.Sprintf("err contains test error=%v", errors.Is(err, errTest)))
	rec.check(t, []string{
		"fn called=false",
		"err contains test error=true",
	})
}
