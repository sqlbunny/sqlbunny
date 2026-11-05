package operations

import (
	"fmt"

	"github.com/sqlbunny/sqlbunny/sqlschema/schema"
)

// DropStatistics drops a Postgres extended statistics object.
type DropStatistics struct {
	SchemaName     string
	TableName      string
	StatisticsName string
}

func (o DropStatistics) GetSQL(transaction bool) string {
	return fmt.Sprintf("DROP STATISTICS %s", sqlName(o.SchemaName, o.StatisticsName))
}

func (o DropStatistics) Apply(d *schema.Database) error {
	s, ok := d.Schemas[o.SchemaName]
	if !ok {
		return fmt.Errorf("no such schema: %s", o.SchemaName)
	}
	t, ok := s.Tables[o.TableName]
	if !ok {
		return fmt.Errorf("no such table: %s", o.TableName)
	}
	if _, ok := t.Statistics[o.StatisticsName]; !ok {
		return fmt.Errorf("no such statistics: %s", o.StatisticsName)
	}
	delete(t.Statistics, o.StatisticsName)
	return nil
}
