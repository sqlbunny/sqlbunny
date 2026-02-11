package operations

import (
	"bytes"
	"fmt"

	"github.com/sqlbunny/sqlbunny/sqlschema/schema"
)

type AlterTable struct {
	SchemaName string
	TableName  string
	Ops        []AlterTableSuboperation
}

func (o AlterTable) GetSQL() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("ALTER TABLE %s\n", sqlName(o.SchemaName, o.TableName)))
	first := true
	for _, op := range o.Ops {
		if !first {
			buf.WriteString(",\n")
		}
		buf.WriteString("    ")
		buf.WriteString(op.GetAlterTableSQL(&o))
		first = false
	}
	return buf.String()
}

func (o AlterTable) Apply(d *schema.Database) error {
	s, ok := d.Schemas[o.SchemaName]
	if !ok {
		return fmt.Errorf("no such schema: %s", o.SchemaName)
	}
	t, ok := s.Tables[o.TableName]
	if !ok {
		return fmt.Errorf("no such table: %s", o.TableName)
	}
	for _, op := range o.Ops {
		err := op.Apply(d, t, o)
		if err != nil {
			return fmt.Errorf("%T on table %s: %w", op, o.TableName, err)
		}
	}
	return nil
}
