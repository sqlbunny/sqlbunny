package operations

import (
	"fmt"
	"io"

	"github.com/sqlbunny/sqlschema/schema"
)

type DropTable struct {
	Name string
}

func (o DropTable) GetSQL() string {
	return fmt.Sprintf("DROP TABLE \"%s\"", o.Name)
}

func (o DropTable) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.DropTable {\n")
	fmt.Fprint(w, "Name: "+esc(o.Name)+",\n")
	fmt.Fprint(w, "}")
}

func (o DropTable) Apply(s *schema.Schema) error {
	if _, ok := s.Tables[o.Name]; !ok {
		return fmt.Errorf("no such table: %s", o.Name)
	}
	delete(s.Tables, o.Name)
	return nil
}
