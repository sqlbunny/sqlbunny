package operations

import (
	"bytes"
	"fmt"

	"github.com/sqlbunny/sqlschema/schema"
)

type DropTable struct {
	Name string
}

func (o DropTable) GetSQL() string {
	return fmt.Sprintf("DROP TABLE \"%s\"", o.Name)
}

func (o DropTable) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.DropTable {\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("}")
}

func (o DropTable) Apply(s *schema.Schema) error {
	if _, ok := s.Tables[o.Name]; !ok {
		return fmt.Errorf("DropTable on non-existing table: %s", o.Name)
	}
	delete(s.Tables, o.Name)
	return nil
}
