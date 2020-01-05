package operations

import (
	"fmt"
	"io"

	"github.com/sqlbunny/sqlschema/schema"
)

type DropTable struct {
	TableName string
}

func (o DropTable) GetSQL() string {
	return fmt.Sprintf("DROP TABLE \"%s\"", o.TableName)
}

func (o DropTable) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.DropTable {\n")
	fmt.Fprint(w, "TableName: "+esc(o.TableName)+",\n")
	fmt.Fprint(w, "}")
}

func (o DropTable) Apply(s *schema.Schema) error {
	if _, ok := s.Tables[o.TableName]; !ok {
		return fmt.Errorf("no such table: %s", o.TableName)
	}
	delete(s.Tables, o.TableName)
	return nil
}
