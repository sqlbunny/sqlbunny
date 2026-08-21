import (
	"context"
	"fmt"
	"strings"

	"github.com/sqlbunny/sqlbunny/runtime/queries"
	"github.com/sqlbunny/sqlbunny/runtime/qm"
)

// NewQuery initializes a new Query using the passed in QueryMods
func NewQuery(mods ...qm.QueryMod) *queries.Query {
	q := &queries.Query{}
	qm.Apply(q, mods...)

	return q
}

// paramList returns "$1, $2, ..., $n" for one row of n values.
func paramList(n int) string {
	parts := make([]string, n)
	for i := range parts {
		parts[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(parts, ", ")
}
