package queries

import (
	"reflect"
	"strings"
	"testing"
)

func TestExpandExprs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		clause     string
		args       []any
		wantClause string
		wantArgs   []any
	}{
		{
			name:       "no exprs is a passthrough",
			clause:     "a=? AND b=?",
			args:       []any{1, 2},
			wantClause: "a=? AND b=?",
			wantArgs:   []any{1, 2},
		},
		{
			name:       "no args is a passthrough",
			clause:     "a IS NOT NULL",
			args:       nil,
			wantClause: "a IS NOT NULL",
			wantArgs:   nil,
		},
		{
			name:       "single expr expands wrapped in parens",
			clause:     "foo = ?",
			args:       []any{NewExpr("? + ?", 25, 34)},
			wantClause: "foo = (? + ?)",
			wantArgs:   []any{25, 34},
		},
		{
			name:       "nested exprs flatten in order",
			clause:     "foo = ?",
			args:       []any{NewExpr("? * ?", NewExpr("? + ?", 25, 34), 45)},
			wantClause: "foo = ((? + ?) * ?)",
			wantArgs:   []any{25, 34, 45},
		},
		{
			name:       "expr mixed with plain args keeps order",
			clause:     "a=? AND b=? AND c=?",
			args:       []any{1, NewExpr("? + ?", 2, 3), 4},
			wantClause: "a=? AND b=(? + ?) AND c=?",
			wantArgs:   []any{1, 2, 3, 4},
		},
		{
			name:       "multiple exprs",
			clause:     "? < x AND x < ?",
			args:       []any{NewExpr("? - ?", 10, 1), NewExpr("? + ?", 20, 1)},
			wantClause: "(? - ?) < x AND x < (? + ?)",
			wantArgs:   []any{10, 1, 20, 1},
		},
		{
			name:       "escaped placeholder consumes no arg",
			clause:     `path \? a = ?`,
			args:       []any{NewExpr("? + ?", 1, 2)},
			wantClause: `path \? a = (? + ?)`,
			wantArgs:   []any{1, 2},
		},
		{
			name:       "expr with no placeholders or args",
			clause:     "x = ?",
			args:       []any{NewExpr("now()")},
			wantClause: "x = (now())",
			wantArgs:   []any{},
		},
		{
			name:       "expr whose arg is a plain slice stays a single arg",
			clause:     "x = ?",
			args:       []any{NewExpr("array[?]", []int{1, 2, 3})},
			wantClause: "x = (array[?])",
			wantArgs:   []any{[]int{1, 2, 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClause, gotArgs := expandExprs(tt.clause, tt.args)
			if gotClause != tt.wantClause {
				t.Errorf("clause:\n got %q\nwant %q", gotClause, tt.wantClause)
			}
			// Treat nil and empty arg slices as equivalent.
			if len(gotArgs) != len(tt.wantArgs) || (len(gotArgs) > 0 && !reflect.DeepEqual(gotArgs, tt.wantArgs)) {
				t.Errorf("args:\n got %#v\nwant %#v", gotArgs, tt.wantArgs)
			}
		})
	}
}

// TestExprPlaceholderNumbering checks that an Expr composed into a WHERE clause
// ends up with correctly numbered ($N) placeholders in the final SQL, proving
// the expansion integrates with the regular placeholder pipeline.
func TestExprPlaceholderNumbering(t *testing.T) {
	t.Parallel()

	a := NewExpr("? + ?", 25, 34)
	b := NewExpr("? * ?", a, 45)

	q := &Query{
		dialect: &Dialect{LQ: '"', RQ: '"', IndexPlaceholders: true},
		from:    []string{"t"},
	}
	AppendWhere(q, "foo = ?", b)
	AppendWhere(q, "bar = ?", 99)

	sql, args := buildQuery(q)

	const wantSQL = `SELECT * FROM "t" WHERE (foo = (($1 + $2) * $3)) AND (bar = $4);`
	if sql != wantSQL {
		t.Errorf("sql:\n got %q\nwant %q", sql, wantSQL)
	}
	wantArgs := []any{25, 34, 45, 99}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Errorf("args:\n got %#v\nwant %#v", args, wantArgs)
	}
}

// TestExprNonIndexedPlaceholders checks the non-Postgres path where `?` are
// kept verbatim instead of being numbered.
func TestExprNonIndexedPlaceholders(t *testing.T) {
	t.Parallel()

	q := &Query{
		dialect: &Dialect{LQ: '`', RQ: '`', IndexPlaceholders: false},
		from:    []string{"t"},
	}
	AppendWhere(q, "foo = ?", NewExpr("? + ?", 1, 2))

	sql, args := buildQuery(q)

	const wantSQL = "SELECT * FROM `t` WHERE (foo = (? + ?));"
	if sql != wantSQL {
		t.Errorf("sql:\n got %q\nwant %q", sql, wantSQL)
	}
	wantArgs := []any{1, 2}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Errorf("args:\n got %#v\nwant %#v", args, wantArgs)
	}
}

// TestExprPanicsOnMismatch checks that a `?`/arg count mismatch panics instead
// of silently letting garbage through.
func TestExprPanicsOnMismatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func()
		expect string
	}{
		{
			name:   "too few args",
			fn:     func() { NewExpr("? + ?", 1) },
			expect: "more ? placeholders than args",
		},
		{
			name:   "too many args",
			fn:     func() { NewExpr("? + ?", 1, 2, 3) },
			expect: "more args than ? placeholders",
		},
		{
			name:   "no placeholders but args given",
			fn:     func() { NewExpr("now()", 1) },
			expect: "more args than ? placeholders",
		},
		{
			name:   "mismatch surfaces through nested expr at construction",
			fn:     func() { NewExpr("? * ?", NewExpr("? + ?", 1), 2) },
			expect: "more ? placeholders than args",
		},
		{
			name:   "mismatch when expanded into a where clause",
			fn:     func() { AppendWhere(&Query{}, "a=? AND b=?", NewExpr("? + ?", 1, 2)) },
			expect: "more ? placeholders than args",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatalf("expected panic, got none")
				}
				msg, _ := r.(string)
				if !strings.Contains(msg, tt.expect) {
					t.Fatalf("panic %q does not contain %q", msg, tt.expect)
				}
			}()
			tt.fn()
		})
	}
}

// TestExprNoExprArgsUnchecked confirms that plain clauses with no Expr args are
// left untouched and unchecked, so `$N`-style placeholders (args but no `?`)
// keep working.
func TestExprNoExprArgsUnchecked(t *testing.T) {
	t.Parallel()

	q := &Query{}
	// 2 args, zero `?` — would trip the arg-count check if validated.
	AppendInnerJoin(q, "thing=$1 AND stuff=$2", 2, 5)
	if q.joins[0].clause != "thing=$1 AND stuff=$2" {
		t.Errorf("clause was rewritten: %q", q.joins[0].clause)
	}
}
