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
func Select(fields ...string) QueryMod {
	return func(q *queries.Query) {
		queries.AppendSelect(q, fields...)
	}
}

// Where allows you to specify a where clause for your statement
func Where(clause string, args ...any) QueryMod {
	return func(q *queries.Query) {
		queries.AppendWhere(q, clause, args...)
	}
}

// WhereIn allows you to specify a "x IN (set)" clause for your where statement
// Example clauses: "field in ?", "(field1,field2) in ?"
func WhereIn(clause string, args ...any) QueryMod {
	return func(q *queries.Query) {
		queries.AppendIn(q, clause, args...)
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
