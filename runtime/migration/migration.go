package migration

import (
	"context"
	"fmt"

	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"github.com/sqlbunny/sqlbunny/sqlschema/operations"
)

type Migration struct {
	Name         string
	Dependencies []string
	Operations   []operations.Operation

	// Transaction, when true, runs all the migration's operations (and the insert
	// that marks the migration as applied) inside a single database transaction,
	// so the migration either fully applies or not at all. Leave false for
	// migrations that must run outside a transaction (e.g. CREATE INDEX
	// CONCURRENTLY, or very large data backfills you don't want to hold a single
	// transaction open for).
	Transaction bool
}

func (m Migration) Run(ctx context.Context) error {
	for _, op := range m.Operations {
		sql := op.GetSQL(m.Transaction)

		_, err := bunny.Exec(ctx, sql)
		if err != nil {
			return fmt.Errorf("error applying operation %v: %w", op, err)
		}
	}
	return nil
}
