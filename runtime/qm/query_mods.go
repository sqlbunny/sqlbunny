package qm

import "github.com/sqlbunny/sqlbunny/runtime/queries"

// QueryMod to modify the query object
type QueryMod func(q *queries.Query)

// Apply the query mods to the Query object
func Apply(q *queries.Query, mods ...QueryMod) {
	for _, mod := range mods {
		mod(q)
	}
}

// SQL allows you to execute a plain SQL statement
func SQL(sql string, args ...any) QueryMod {
	return func(q *queries.Query) {
		queries.SetSQL(q, sql, args...)
	}
}

// Load allows you to specify foreign key relationships to eager load
// for your query. Passed in relationships need to be in the format
// MyThing or MyThings.
// Relationship name plurality is important, if your relationship is
// singular, you need to specify the singular form and vice versa.
func Load(relationships ...string) QueryMod {
	return func(q *queries.Query) {
		queries.AppendLoad(q, relationships...)
	}
}

// InnerJoin on another model
func InnerJoin(clause string, args ...any) QueryMod {
	return func(q *queries.Query) {
		queries.AppendInnerJoin(q, clause, args...)
	}
}

// Select specific fields opposed to all fields
func Select[T ~string](fields ...T) QueryMod {
	s := make([]string, len(fields))
	for i, f := range fields {
		s[i] = string(f)
	}
	return func(q *queries.Query) {
		queries.AppendSelect(q, s...)
	}
}

// Expr builds a composable SQL expression from a clause containing `?`
// placeholders and its args. The returned value can be passed as an arg to
// Where (or any other clause taking `?` placeholders), where it expands in
// place — wrapped in parentheses — with its args spliced into the surrounding
// arg list. Exprs nest: an Expr's args may themselves be Exprs.
//
//	a := qm.Expr("? + ?", 25, 34)
//	b := qm.Expr("? * ?", a, 45)
//	qm.Where("foo = ?", b)
//
// is equivalent to:
//
//	qm.Where("foo = ((? + ?) * ?)", 25, 34, 45)
func Expr(clause string, args ...any) *queries.Expr {
	return queries.NewExpr(clause, args...)
}

// Where allows you to specify a where clause for your statement
func Where(clause string, args ...any) QueryMod {
	return func(q *queries.Query) {
		queries.AppendWhere(q, clause, args...)
	}
}

// GroupBy allows you to specify a group by clause for your statement
func GroupBy(clause string) QueryMod {
	return func(q *queries.Query) {
		queries.AppendGroupBy(q, clause)
	}
}

// OrderBy allows you to specify a order by clause for your statement
func OrderBy(clause string) QueryMod {
	return func(q *queries.Query) {
		queries.AppendOrderBy(q, clause)
	}
}

// Having allows you to specify a having clause for your statement
func Having(clause string, args ...any) QueryMod {
	return func(q *queries.Query) {
		queries.AppendHaving(q, clause, args...)
	}
}

// From allows to specify the model for your statement
func From(from string) QueryMod {
	return func(q *queries.Query) {
		queries.AppendFrom(q, from)
	}
}

// Limit the number of returned rows
func Limit(limit int) QueryMod {
	return func(q *queries.Query) {
		queries.SetLimit(q, limit)
	}
}

// Offset into the results
func Offset(offset int) QueryMod {
	return func(q *queries.Query) {
		queries.SetOffset(q, offset)
	}
}

// For inserts a concurrency locking clause at the end of your statement
func For(clause string) QueryMod {
	return func(q *queries.Query) {
		queries.SetFor(q, clause)
	}
}
