package operations

import (
	"fmt"
	"io"

	"github.com/sqlbunny/sqlschema/schema"
)

type RenameTable struct {
	SchemaName   string
	TableName    string
	NewTableName string
}

func (o RenameTable) GetSQL() string {
	return fmt.Sprintf("ALTER TABLE %s RENAME TO \"%s\"", sqlName(o.SchemaName, o.TableName), o.NewTableName)
}

func (o RenameTable) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.RenameTable {\n")
	fmt.Fprint(w, "SchemaName: "+esc(o.SchemaName)+",\n")
	fmt.Fprint(w, "TableName: "+esc(o.TableName)+",\n")
	fmt.Fprint(w, "NewTableName: "+esc(o.NewTableName)+",\n")
	fmt.Fprint(w, "}")
}

func (o RenameTable) Apply(d *schema.Database) error {
	s, ok := d.Schemas[o.SchemaName]
	if !ok {
		return fmt.Errorf("no such schema: %s", o.SchemaName)
	}
	t, ok := s.Tables[o.TableName]
	if !ok {
		return fmt.Errorf("no such table: %s", o.TableName)
	}
	if _, ok := s.Tables[o.NewTableName]; ok {
		return fmt.Errorf("destination table already exists: %s", o.NewTableName)
	}

	delete(s.Tables, o.TableName)
	s.Tables[o.NewTableName] = t

	for _, t2 := range s.Tables {
		for _, fk := range t2.ForeignKeys {
			if fk.ForeignSchema == o.SchemaName && fk.ForeignTable == o.TableName {
				fk.ForeignTable = o.NewTableName
			}
		}
	}
	return nil
}
