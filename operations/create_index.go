package operations

import (
	"bytes"
	"fmt"

	"github.com/sqlbunny/sqlschema/schema"
)

type CreateIndex struct {
	Name      string
	IndexName string
	Columns   []string
}

func (o CreateIndex) GetSQL() string {
	return fmt.Sprintf("CREATE INDEX CONCURRENTLY \"%s\" ON \"%s\" (%s)", o.IndexName, o.Name, columnList(o.Columns))
}

func (o CreateIndex) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.CreateIndex {\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("IndexName: " + esc(o.IndexName) + ",\n")
	buf.WriteString("Columns: []string{" + columnList(o.Columns) + "},\n")
	buf.WriteString("}")
}

func (o CreateIndex) Apply(s *schema.Schema) error {
	t, ok := s.Tables[o.Name]
	if !ok {
		return fmt.Errorf("CreateIndex on non-existing table: %s", o.Name)
	}
	if _, ok := t.Indexes[o.IndexName]; ok {
		return fmt.Errorf("CreateIndex index already exists: table %s, index %s ", o.Name, o.IndexName)
	}
	t.Indexes[o.IndexName] = &schema.Index{
		Columns: o.Columns,
	}
	return nil
}
