package operations

import (
	"bytes"
	"fmt"

	"github.com/sqlbunny/sqlschema/schema"
)

type AlterTable struct {
	Name string
	Ops  []AlterTableSuboperation
}

func (o AlterTable) GetSQL() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("ALTER TABLE \"%s\"\n", o.Name))
	first := true
	for _, op := range o.Ops {
		if !first {
			buf.WriteString(",\n")
		}
		buf.WriteString("    ")
		buf.WriteString(op.GetAlterTableSQL(&o))
		first = false
	}
	return buf.String()
}

func (o AlterTable) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.AlterTable {\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("Ops: []operations.AlterTableSuboperation{\n")
	for _, op := range o.Ops {
		op.Dump(buf)
		buf.WriteString(",\n")
	}
	buf.WriteString("},\n")
	buf.WriteString("}")
}

func (o AlterTable) Apply(s *schema.Schema) error {
	t, ok := s.Tables[o.Name]
	if !ok {
		return fmt.Errorf("AlterTable on non-existing table: %s", o.Name)
	}
	for _, op := range o.Ops {
		err := op.Apply(s, t, o)
		if err != nil {
			return fmt.Errorf("AlterTable on table %s: %w", o.Name, err)
		}
	}
	return nil
}
