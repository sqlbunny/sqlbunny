package operations

import (
	"fmt"
	"io"

	"github.com/sqlbunny/sqlschema/schema"
)

type DropSchema struct {
	SchemaName string
}

func (o DropSchema) GetSQL() string {
	return fmt.Sprintf("DROP SCHEMA \"%s\"", o.SchemaName)
}

func (o DropSchema) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.DropSchema {\n")
	fmt.Fprint(w, "SchemaName: "+esc(o.SchemaName)+",\n")
	fmt.Fprint(w, "}")
}

func (o DropSchema) Apply(d *schema.Database) error {
	if _, ok := d.Schemas[o.SchemaName]; !ok {
		return fmt.Errorf("no such schema: %s", o.SchemaName)
	}
	delete(d.Schemas, o.SchemaName)
	return nil
}
