package migration

import (
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
			return fmt.Errorf("error applying operation %v: %w", op, err)
		}
	}
	return nil
}
