package operations

import (
	"fmt"
	"io"

	"github.com/sqlbunny/sqlschema/schema"
)

type DropIndex struct {
	SchemaName string
	TableName  string
	IndexName  string
}

func (o DropIndex) GetSQL() string {
	return fmt.Sprintf("DROP INDEX %s", sqlName(o.SchemaName, o.IndexName))
}

func (o DropIndex) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.DropIndex {\n")
	fmt.Fprint(w, "SchemaName: "+esc(o.SchemaName)+",\n")
	fmt.Fprint(w, "TableName: "+esc(o.TableName)+",\n")
	fmt.Fprint(w, "IndexName: "+esc(o.IndexName)+",\n")
	fmt.Fprint(w, "}")
}

func (o DropIndex) Apply(d *schema.Database) error {
	s, ok := d.Schemas[o.SchemaName]
	if !ok {
		return fmt.Errorf("no such schema: %s", o.SchemaName)
	}
	t, ok := s.Tables[o.TableName]
	if !ok {
		return fmt.Errorf("no such table: %s", o.TableName)
	}
	if _, ok := t.Indexes[o.IndexName]; !ok {
		return fmt.Errorf("no such index: %s", o.IndexName)
	}
	delete(t.Indexes, o.IndexName)
	return nil
}
