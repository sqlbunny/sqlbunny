package operations

import (
	"fmt"
	"io"

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

func (o CreateIndex) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.CreateIndex {\n")
	fmt.Fprint(w, "Name: "+esc(o.Name)+",\n")
	fmt.Fprint(w, "IndexName: "+esc(o.IndexName)+",\n")
	fmt.Fprint(w, "Columns: []string{"+columnList(o.Columns)+"},\n")
	fmt.Fprint(w, "}")
}

func (o CreateIndex) Apply(s *schema.Schema) error {
	t, ok := s.Tables[o.Name]
	if !ok {
		return fmt.Errorf("no such table: %s", o.Name)
	}
	if _, ok := t.Indexes[o.IndexName]; ok {
		return fmt.Errorf("index already exists: %s ", o.IndexName)
	}
	t.Indexes[o.IndexName] = &schema.Index{
		Columns: o.Columns,
	}
	return nil
}
