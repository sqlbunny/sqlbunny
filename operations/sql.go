package operations

import (
	"bytes"

	"github.com/sqlbunny/sqlschema/schema"
)

type SQL struct {
	SQL string
}

func (o SQL) GetSQL() string {
	return o.SQL
}

func (o SQL) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.SQL {\n")
	buf.WriteString("SQL: " + esc(o.SQL) + ",\n")
	buf.WriteString("}")
}

func (o SQL) Apply(s *schema.Schema) error {
	// do nothing.
	return nil
}
