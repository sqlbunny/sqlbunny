package operations

import (
	"fmt"

	"github.com/sqlbunny/sqlschema/schema"
)

type RenameColumn struct {
	SchemaName    string
	TableName     string
	OldColumnName string
	NewColumnName string
}

func (o RenameColumn) GetSQL() string {
	return fmt.Sprintf("ALTER TABLE %s RENAME COLUMN \"%s\" TO \"%s\"", sqlName(o.SchemaName, o.TableName), o.OldColumnName, o.NewColumnName)
}

func (o RenameColumn) Apply(d *schema.Database) error {
	s, ok := d.Schemas[o.SchemaName]
	if !ok {
		return fmt.Errorf("no such schema: %s", o.SchemaName)
	}
	t, ok := s.Tables[o.TableName]
	if !ok {
		return fmt.Errorf("no such table: %s", o.TableName)
	}

	c, ok := t.Columns[o.OldColumnName]
	if !ok {
		return fmt.Errorf("no such column on table %s: %s", o.TableName, o.OldColumnName)
	}

	delete(t.Columns, o.OldColumnName)
	t.Columns[o.NewColumnName] = c

	if t.PrimaryKey != nil {
		for i := range t.PrimaryKey.Columns {
			if t.PrimaryKey.Columns[i] == o.OldColumnName {
				t.PrimaryKey.Columns[i] = o.NewColumnName
			}
		}
	}
	for _, idx := range t.Indexes {
		for i := range idx.Columns {
			if idx.Columns[i] == o.OldColumnName {
				idx.Columns[i] = o.NewColumnName
			}
		}
	}
	for _, idx := range t.Uniques {
		for i := range idx.Columns {
			if idx.Columns[i] == o.OldColumnName {
				idx.Columns[i] = o.NewColumnName
			}
		}
	}

	for _, fk := range t.ForeignKeys {
		for i := range fk.LocalColumns {
			if fk.LocalColumns[i] == o.OldColumnName {
				fk.LocalColumns[i] = o.NewColumnName
			}
		}
	}

	for _, m2 := range s.Tables {
		for _, fk := range m2.ForeignKeys {
			if fk.ForeignTable == o.TableName {
				for _, fk := range m2.ForeignKeys {
					for i := range fk.ForeignColumns {
						if fk.ForeignColumns[i] == o.OldColumnName {
							fk.ForeignColumns[i] = o.NewColumnName
						}
					}
				}
			}
		}
	}
	return nil
}
