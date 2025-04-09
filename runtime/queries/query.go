package queries

import (
	"context"
	"database/sql"

	"github.com/sqlbunny/sqlbunny/runtime/bunny"
)

// joinKind is the type of join
type joinKind int

// Join type constants
const (
	JoinInner joinKind = iota
	JoinOuterLeft
	JoinOuterRight
	JoinNatural
)

// Query holds the state for the built up query
type Query struct {
	dialect    *Dialect
	rawSQL     rawSQL
	load       []string
	delete     bool
	update     map[string]any
	selectCols []string
	count      bool
	from       []string
	joins      []join
	where      []where
	in         []in
	groupBy    []string
	orderBy    []string
	having     []having
	limit      int
	offset     int
	forlock    string
}

// Dialect holds values that direct the query builder
// how to build compatible queries for each database.
// Each database driver needs to implement functions
// that provide these values.
type Dialect struct {
	// The left quote character for SQL identifiers
	LQ byte
	// The right quote character for SQL identifiers
	RQ byte
	// Bool flag indicating whether indexed
	// placeholders ($1) are used, or ? placeholders.
	IndexPlaceholders bool
	// Bool flag indicating whether "TOP" or "LIMIT" clause
	// must be used for rows limitation
	UseTopClause bool
}

type where struct {
	clause string
	args   []any
}

type in struct {
	clause string
	args   []any
}

type having struct {
	clause string
	args   []any
}

type rawSQL struct {
	sql  string
	args []any
}

type join struct {
	kind   joinKind
	clause string
	args   []any
}

// Raw makes a raw query, usually for use with bind
func Raw(query string, args ...any) *Query {
	return &Query{
		rawSQL: rawSQL{
			sql:  query,
			args: args,
		},
	}
}

// Exec executes a query that does not need a row returned
func (q *Query) Exec(ctx context.Context) (sql.Result, error) {
	qs, args := buildQuery(q)
	return bunny.Exec(ctx, qs, args...)
}

// QueryRow executes the query for the One finisher and returns a row
func (q *Query) QueryRow(ctx context.Context) *sql.Row {
	qs, args := buildQuery(q)
	return bunny.QueryRow(ctx, qs, args...)
}

// Query executes the query for the All finisher and returns multiple rows
func (q *Query) Query(ctx context.Context) (*sql.Rows, error) {
	qs, args := buildQuery(q)
	return bunny.Query(ctx, qs, args...)
}

// SetDialect on the query.
func SetDialect(q *Query, dialect *Dialect) {
	q.dialect = dialect
}

// SetSQL on the query.
func SetSQL(q *Query, sql string, args ...any) {
	q.rawSQL = rawSQL{sql: sql, args: args}
}

// SetLoad on the query.
func SetLoad(q *Query, relationships ...string) {
	q.load = append([]string(nil), relationships...)
}

// AppendLoad on the query.
func AppendLoad(q *Query, relationships ...string) {
	q.load = append(q.load, relationships...)
}

// SetSelect on the query.
func SetSelect(q *Query, sel []string) {
	q.selectCols = sel
}

// GetSelect from the query
func GetSelect(q *Query) []string {
	return q.selectCols
}

// SetCount on the query.
func SetCount(q *Query) {
	q.count = true
}

// SetDelete on the query.
func SetDelete(q *Query) {
	q.delete = true
}

// SetLimit on the query.
func SetLimit(q *Query, limit int) {
	q.limit = limit
}

// SetOffset on the query.
func SetOffset(q *Query, offset int) {
	q.offset = offset
}

// SetFor on the query.
func SetFor(q *Query, clause string) {
	q.forlock = clause
}

// SetUpdate on the query.
func SetUpdate(q *Query, cols map[string]any) {
	q.update = cols
}

// AppendSelect on the query.
func AppendSelect(q *Query, fields ...string) {
	q.selectCols = append(q.selectCols, fields...)
}

// AppendFrom on the query.
func AppendFrom(q *Query, from ...string) {
	q.from = append(q.from, from...)
}

// SetFrom replaces the current from statements.
func SetFrom(q *Query, from ...string) {
	q.from = append([]string(nil), from...)
}

// AppendInnerJoin on the query.
func AppendInnerJoin(q *Query, clause string, args ...any) {
	q.joins = append(q.joins, join{clause: clause, kind: JoinInner, args: args})
}

// AppendHaving on the query.
func AppendHaving(q *Query, clause string, args ...any) {
	q.having = append(q.having, having{clause: clause, args: args})
}

// AppendWhere on the query.
func AppendWhere(q *Query, clause string, args ...any) {
	q.where = append(q.where, where{clause: clause, args: args})
}

// AppendIn on the query.
func AppendIn(q *Query, clause string, args ...any) {
	q.in = append(q.in, in{clause: clause, args: args})
}

// AppendGroupBy on the query.
func AppendGroupBy(q *Query, clause string) {
	q.groupBy = append(q.groupBy, clause)
}

// AppendOrderBy on the query.
func AppendOrderBy(q *Query, clause string) {
	q.orderBy = append(q.orderBy, clause)
}
