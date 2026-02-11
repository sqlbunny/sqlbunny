package operations

import (
	"fmt"

	"github.com/sqlbunny/sqlschema/schema"
)

type DropTable struct {
	SchemaName string
	TableName  string
}

func (o DropTable) GetSQL() string {
	return fmt.Sprintf("DROP TABLE %s", sqlName(o.SchemaName, o.TableName))
}
func (o DropTable) Apply(d *schema.Database) error {
	s, ok := d.Schemas[o.SchemaName]
	if !ok {
		return fmt.Errorf("no such schema: %s", o.SchemaName)
	}
	if _, ok := s.Tables[o.TableName]; !ok {
		return fmt.Errorf("no such table: %s", o.TableName)
	}
	delete(s.Tables, o.TableName)
	return nil
}
