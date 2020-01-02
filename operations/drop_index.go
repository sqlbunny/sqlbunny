package operations

import (
	"fmt"
	"io"

	"github.com/sqlbunny/sqlschema/schema"
)

type DropIndex struct {
	Name      string
	IndexName string
}

func (o DropIndex) GetSQL() string {
	return fmt.Sprintf("DROP INDEX \"%s\"", o.IndexName)
}

func (o DropIndex) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.DropIndex {\n")
	fmt.Fprint(w, "Name: "+esc(o.Name)+",\n")
	fmt.Fprint(w, "IndexName: "+esc(o.IndexName)+",\n")
	fmt.Fprint(w, "}")
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
