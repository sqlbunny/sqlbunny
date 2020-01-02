package operations

import (
	"bytes"
	"fmt"

	"github.com/sqlbunny/sqlschema/schema"
)

type RenameTable struct {
	OldName string
	NewName string
}

func (o RenameTable) GetSQL() string {
	return fmt.Sprintf("ALTER TABLE \"%s\" RENAME TO \"%s\"", o.OldName, o.NewName)
}

func (o RenameTable) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.RenameTable {\n")
	buf.WriteString("OldName: " + esc(o.OldName) + ",\n")
	buf.WriteString("NewName: " + esc(o.NewName) + ",\n")
	buf.WriteString("}")
}

func (o RenameTable) Apply(s *schema.Schema) error {
	t, ok := s.Tables[o.OldName]
	if !ok {
		return fmt.Errorf("RenameTable on non-existing table: %s", o.OldName)
	}
	if _, ok := s.Tables[o.NewName]; ok {
		return fmt.Errorf("RenameTable new table name already exists: %s", o.NewName)
	}

	delete(s.Tables, o.OldName)
	s.Tables[o.NewName] = t

	for _, t2 := range s.Tables {
		for _, fk := range t2.ForeignKeys {
			if fk.ForeignTable == o.OldName {
				fk.ForeignTable = o.NewName
			}
		}
	}
	return nil
}
