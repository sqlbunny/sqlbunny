import (
	"context"
	"fmt"
	"strings"

	"github.com/sqlbunny/sqlbunny/runtime/queries"
	"github.com/sqlbunny/sqlbunny/runtime/qm"
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

// paramList returns "$1, $2, ..., $n" (or "?, ?, ...") for one row of n values.
func paramList(n int) string {
	parts := make([]string, n)
	for i := range parts {
		if dialect.IndexPlaceholders {
			parts[i] = fmt.Sprintf("$%d", i+1)
		} else {
			parts[i] = "?"
		}
	}
	return strings.Join(parts, ", ")
}
