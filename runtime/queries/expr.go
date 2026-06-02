package queries

import (
	"fmt"
	"strings"
)

// Expr is a composable SQL expression: a clause string containing `?`
// placeholders together with the args that fill them. An Expr can be passed
// as an arg wherever a `?` placeholder is expected, and it will be expanded
// in place (wrapped in parentheses) with its args spliced into the
// surrounding arg list.
//
// For example:
//
//	a := qm.Expr("? + ?", 25, 34)
//	b := qm.Expr("? * ?", a, 45)
//	qm.Where("foo = ?", b)
//
// is equivalent to:
//
//	qm.Where("foo = ((? + ?) * ?)", 25, 34, 45)
type Expr struct {
	clause string
	args   []any
}

// NewExpr builds an Expr from a clause and its args. It expands eagerly: any
// Expr args are spliced in immediately and the clause/args are validated now,
// so a malformed Expr panics at the construction site rather than later when
// the query is built. The stored Expr is therefore always fully flattened —
// its clause contains only literal `?` placeholders and its args contain no
// Exprs.
func NewExpr(clause string, args ...any) *Expr {
	clause, args = expandExprs(clause, args)
	return &Expr{clause: clause, args: args}
}

// containsExpr reports whether any arg is an *Expr.
func containsExpr(args []any) bool {
	for _, a := range args {
		if _, ok := a.(*Expr); ok {
			return true
		}
	}
	return false
}

// expandExprs walks clause, and for every unescaped `?` placeholder consumes
// the next arg. If that arg is an *Expr, its (already fully expanded) clause is
// spliced in wrapped in parentheses and its args are inlined; otherwise the `?`
// and the arg are kept as-is. The result is a clause whose only remaining `?`
// are literal placeholders, each backed by exactly one non-Expr arg, in order —
// ready for the regular placeholder-numbering pipeline.
//
// Escaped placeholders (`\?`) are passed through untouched and do not consume
// an arg.
//
// The clause and its args must line up exactly: a `?`/arg count mismatch is a
// programmer error and panics. Because NewExpr expands eagerly, stored Exprs are
// already flat, so no recursion is needed here.
//
// Callers should only invoke this when an arg is actually an Expr (see
// containsExpr): plain clauses may use `$N`-style placeholders, which carry args
// but no `?` and would otherwise trip the arg-count check.
func expandExprs(clause string, args []any) (string, []any) {
	var sb strings.Builder
	outArgs := make([]any, 0, len(args))
	argIdx := 0

	for i := 0; i < len(clause); i++ {
		c := clause[i]
		if c == '\\' && i+1 < len(clause) && clause[i+1] == '?' {
			// Escaped placeholder: emit verbatim, consume no arg.
			sb.WriteByte('\\')
			sb.WriteByte('?')
			i++
			continue
		}
		if c != '?' {
			sb.WriteByte(c)
			continue
		}

		// Unescaped placeholder.
		if argIdx >= len(args) {
			panic(fmt.Sprintf("expr: more ? placeholders than args in clause %q (%d args)", clause, len(args)))
		}
		arg := args[argIdx]
		argIdx++

		if e, ok := arg.(*Expr); ok {
			// e is already fully expanded by NewExpr; inline it directly.
			sb.WriteByte('(')
			sb.WriteString(e.clause)
			sb.WriteByte(')')
			outArgs = append(outArgs, e.args...)
		} else {
			sb.WriteByte('?')
			outArgs = append(outArgs, arg)
		}
	}

	// Every arg must have been consumed by a placeholder.
	if argIdx != len(args) {
		panic(fmt.Sprintf("expr: more args than ? placeholders in clause %q (%d args, %d consumed)", clause, len(args), argIdx))
	}

	return sb.String(), outArgs
}
