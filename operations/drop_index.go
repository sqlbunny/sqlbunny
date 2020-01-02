package operations

import (
	"bytes"
	"fmt"

	"github.com/sqlbunny/sqlschema/schema"
)

type DropIndex struct {
	Name      string
	IndexName string
}

func (o DropIndex) GetSQL() string {
	return fmt.Sprintf("DROP INDEX \"%s\"", o.IndexName)
}

func (o DropIndex) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.DropIndex {\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("IndexName: " + esc(o.IndexName) + ",\n")
	buf.WriteString("}")
}

func (o DropIndex) Apply(s *schema.Schema) error {
	t, ok := s.Tables[o.Name]
	if !ok {
		return fmt.Errorf("DropIndex on non-existing table: %s", o.Name)
	}
	if _, ok := t.Indexes[o.IndexName]; !ok {
		return fmt.Errorf("DropIndex index doesn't exist: table %s, index %s ", o.Name, o.IndexName)
	}
	delete(t.Indexes, o.IndexName)
	return nil
}
