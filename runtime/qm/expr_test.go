package qm_test

import (
	"testing"

	"github.com/sqlbunny/sqlbunny/runtime/queries"
	"github.com/sqlbunny/sqlbunny/runtime/qm"
)

// TestExprComposition exercises the public qm.Expr API exactly as documented:
// nested Exprs compose, and applying them as query mods doesn't panic. The
// actual SQL expansion is asserted in the queries package tests, which can
// reach the unexported query internals.
func TestExprComposition(t *testing.T) {
	t.Parallel()

	a := qm.Expr("? + ?", 25, 34)
	b := qm.Expr("? * ?", a, 45)
	if b == nil {
		t.Fatal("qm.Expr returned nil")
	}

	var q queries.Query
	qm.Apply(&q,
		qm.From("t"),
		qm.Where("foo = ?", b),
	)
}
