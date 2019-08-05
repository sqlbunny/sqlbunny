package migration

import (
	"bytes"
	"context"
	"fmt"
)

type Migration struct {
	Name         string
	Dependencies []string
	Operations   []Operation
}

func (m Migration) Run(ctx context.Context) error {
	for _, op := range m.Operations {
		err := op.Run(ctx)
		if err != nil {
			return err
		}
	}
	return nil
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
	buf.WriteString("Operations: []migration.Operation{\n")
	for _, op := range m.Operations {
		op.Dump(buf)
		buf.WriteString(",\n")
	}
	buf.WriteString("},\n")
	buf.WriteString("}")
}
