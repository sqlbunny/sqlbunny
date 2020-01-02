package migration

import (
	"context"
	"fmt"
	"io"

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

func (m Migration) Dump(w io.Writer) {
	fmt.Fprint(w, "&migration.Migration {\n")
	fmt.Fprintf(w, "Name: %s,\n", esc(m.Name))
	fmt.Fprint(w, "Dependencies: []string{")
	for i, d := range m.Dependencies {
		if i != 0 {
			fmt.Fprint(w, ",")
		}
		fmt.Fprint(w, esc(d))
	}
	fmt.Fprint(w, "},\n")
	fmt.Fprint(w, "Operations: []operations.Operation{\n")
	for _, op := range m.Operations {
		op.Dump(w)
		fmt.Fprint(w, ",\n")
	}
	fmt.Fprint(w, "},\n")
	fmt.Fprint(w, "}")
}
