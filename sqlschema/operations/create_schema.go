package operations

import (
	"fmt"

	"github.com/sqlbunny/sqlschema/schema"
)

type CreateSchema struct {
	SchemaName string
}

func (o CreateSchema) GetSQL() string {
	return fmt.Sprintf("CREATE SCHEMA \"%s\"", o.SchemaName)
}

func (o CreateSchema) Apply(d *schema.Database) error {
	if _, ok := d.Schemas[o.SchemaName]; ok {
		return fmt.Errorf("schema already exists: %s", o.SchemaName)
	}
	d.Schemas[o.SchemaName] = schema.NewSchema()
	return nil
}
