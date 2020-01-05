package operations

import (
	"bytes"
	"fmt"
	"io"

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

func (o AlterTable) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.AlterTable {\n")
	fmt.Fprint(w, "Name: "+esc(o.Name)+",\n")
	fmt.Fprint(w, "Ops: []operations.AlterTableSuboperation{\n")
	for _, op := range o.Ops {
		op.Dump(w)
		fmt.Fprint(w, ",\n")
	}
	fmt.Fprint(w, "},\n")
	fmt.Fprint(w, "}")
}

func (o AlterTable) Apply(s *schema.Schema) error {
	t, ok := s.Tables[o.Name]
	if !ok {
		return fmt.Errorf("no such table: %s", o.Name)
	}
	for _, op := range o.Ops {
		err := op.Apply(s, t, o)
		if err != nil {
			return fmt.Errorf("%T on table %s: %w", op, o.Name, err)
		}
	}
	return nil
}
