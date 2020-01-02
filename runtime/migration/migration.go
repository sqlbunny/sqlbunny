package migration

import (
	"bytes"
	"context"
	"fmt"

	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"github.com/sqlbunny/sqlschema/operations"
)

type Migration struct {
	Name         string
	Dependencies []string
	Operations   []operations.Operation
}

func (m Migration) Run(ctx context.Context) error {
	for _, op := range m.Operations {
		sql := op.GetSQL()

		_, err := bunny.Exec(ctx, sql)
		if err != nil {
			return err
		}
	}
	return nil
}

func esc(s string) string {
	return fmt.Sprintf("%#v", s)
}

func (m Migration) Dump(buf *bytes.Buffer) {
	buf.WriteString("&migration.Migration {\n")
	buf.WriteString(fmt.Sprintf("Name: %s,\n", esc(m.Name)))
	buf.WriteString("Dependencies: []string{")
	for i, d := range m.Dependencies {
		if i != 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(esc(d))
	}
	buf.WriteString("},\n")
	buf.WriteString("Operations: []operations.Operation{\n")
	for _, op := range m.Operations {
		op.Dump(buf)
		buf.WriteString(",\n")
	}
	buf.WriteString("},\n")
	buf.WriteString("}")
}
