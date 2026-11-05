package operations

import (
	"fmt"
	"slices"
	"strings"

	"github.com/sqlbunny/sqlbunny/sqlschema/schema"
)

// CreateStatistics creates a Postgres extended statistics object
// (CREATE STATISTICS). It tells the planner how columns relate to each other,
// e.g. that one column determines another, so it stops treating conditions
// on them as independent. Postgres fills it on the next ANALYZE.
type CreateStatistics struct {
	SchemaName     string
	TableName      string
	StatisticsName string
	Columns        []string // At least 2.
	Kinds          []string // ndistinct, dependencies, mcv. If empty, Postgres builds all kinds.
}

func (o CreateStatistics) GetSQL(transaction bool) string {
	kinds := ""
	if len(o.Kinds) != 0 {
		kinds = " (" + strings.Join(o.Kinds, ", ") + ")"
	}
	return fmt.Sprintf("CREATE STATISTICS %s%s ON %s FROM %s",
		sqlName(o.SchemaName, o.StatisticsName), kinds, columnList(o.Columns), sqlName(o.SchemaName, o.TableName))
}

func (o CreateStatistics) Apply(d *schema.Database) error {
	s, ok := d.Schemas[o.SchemaName]
	if !ok {
		return fmt.Errorf("no such schema: %s", o.SchemaName)
	}
	t, ok := s.Tables[o.TableName]
	if !ok {
		return fmt.Errorf("no such table: %s", o.TableName)
	}
	if t.Statistics == nil {
		t.Statistics = make(map[string]*schema.Statistics)
	}
	if _, ok := t.Statistics[o.StatisticsName]; ok {
		return fmt.Errorf("statistics already exist: %s", o.StatisticsName)
	}
	t.Statistics[o.StatisticsName] = &schema.Statistics{
		Columns: slices.Clone(o.Columns),
		Kinds:   slices.Clone(o.Kinds),
	}
	return nil
}
