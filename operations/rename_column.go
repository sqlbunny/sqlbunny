package operations

import (
	"bytes"
	"fmt"

	"github.com/sqlbunny/sqlschema/schema"
)

type RenameColumn struct {
	Name          string
	OldColumnName string
	NewColumnName string
}

func (o RenameColumn) GetSQL() string {
	return fmt.Sprintf("ALTER TABLE \"%s\" RENAME COLUMN \"%s\" TO \"%s\"", o.Name, o.OldColumnName, o.NewColumnName)
}

func (o RenameColumn) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.RenameColumn {Name: " + esc(o.Name) + ", OldColumnName: " + esc(o.OldColumnName) + ", NewColumnName: " + esc(o.NewColumnName) + "}")
}

func (o RenameColumn) Apply(s *schema.Schema) error {
	t, ok := s.Tables[o.Name]
	if !ok {
		return fmt.Errorf("RenameColumn on non-existing table: %s", o.Name)
	}

	c, ok := t.Columns[o.OldColumnName]
	if !ok {
		return fmt.Errorf("RenameColumn non-existing: table %s, column %s", o.Name, o.OldColumnName)
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
			if fk.ForeignTable == o.Name {
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
