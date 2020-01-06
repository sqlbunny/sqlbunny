package operations

import (
	"fmt"
	"io"

	"github.com/sqlbunny/sqlschema/schema"
)

type SQL struct {
	SQL string
}

func (o SQL) GetSQL() string {
	return o.SQL
}

func (o SQL) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.SQL {\n")
	fmt.Fprint(w, "SQL: "+esc(o.SQL)+",\n")
	fmt.Fprint(w, "}")
}

func (o SQL) Apply(d *schema.Database) error {
	// do nothing.
	return nil
}
