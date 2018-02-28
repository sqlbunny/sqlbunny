import (
	"context"

	"github.com/KernelPay/sqlboiler/boil/queries"
	"github.com/KernelPay/sqlboiler/boil/qm"
)

var dialect = queries.Dialect{
	LQ: 0x{{printf "%x" .Dialect.LQ}},
	RQ: 0x{{printf "%x" .Dialect.RQ}},
	IndexPlaceholders: {{.Dialect.IndexPlaceholders}},
	UseTopClause: {{.Dialect.UseTopClause}},
}

// NewQuery initializes a new Query using the passed in QueryMods
func NewQuery(ctx context.Context, mods ...qm.QueryMod) *queries.Query {
	q := &queries.Query{}
	queries.SetContext(q, ctx)
	queries.SetDialect(q, &dialect)
	qm.Apply(q, mods...)

	return q
}
