import (
	"context"

	"github.com/kernelpayments/sqlbunny/runtime/queries"
	"github.com/kernelpayments/sqlbunny/runtime/qm"
)

var dialect = queries.Dialect{
	LQ: 0x{{printf "%x" .Dialect.LQ}},
	RQ: 0x{{printf "%x" .Dialect.RQ}},
	IndexPlaceholders: {{.Dialect.IndexPlaceholders}},
	UseTopClause: {{.Dialect.UseTopClause}},
}

// NewQuery initializes a new Query using the passed in QueryMods
func NewQuery(mods ...qm.QueryMod) *queries.Query {
	q := &queries.Query{}
	queries.SetDialect(q, &dialect)
	qm.Apply(q, mods...)

	return q
}
