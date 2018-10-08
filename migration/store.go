package migration

import (
	"context"
	"fmt"
	"time"

	"github.com/KernelPay/sqlbunny/bunny"
	"github.com/KernelPay/sqlbunny/schema"
)

type MigrationStore struct {
	migrations map[int64]OperationList
}

const (
	checkMigrationsTableSQL  = "SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'migrations'"
	createMigrationsTableSQL = "CREATE TABLE migrations (id integer PRIMARY KEY, time timestamptz)"
	insertMigrationSQL       = "INSERT INTO migrations (id, time) VALUES($1, $2)"
	maxMigrationSQL          = "SELECT coalesce(max(id), 0) from migrations"
)

func (m *MigrationStore) Run(ctx context.Context) error {
	var count int64
	if err := bunny.QueryRow(ctx, checkMigrationsTableSQL).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		if _, err := bunny.Exec(ctx, createMigrationsTableSQL); err != nil {
			return err
		}
	}

	var max int64
	if err := bunny.QueryRow(ctx, maxMigrationSQL).Scan(&max); err != nil {
		return err
	}

	var i int64 = 1
	for {
		ops, ok := m.migrations[i]
		if !ok {
			break
		}

		if i > max {
			for _, op := range ops {
				err := op.Run(ctx)
				if err != nil {
					return err
				}
			}

			if _, err := bunny.Exec(ctx, insertMigrationSQL, i, time.Now()); err != nil {
				return err
			}
		}

		i++
	}
	return nil
}

func (m *MigrationStore) Register(index int64, ops OperationList) {
	if m.migrations == nil {
		m.migrations = make(map[int64]OperationList)
	}
	if _, ok := m.migrations[index]; ok {
		panic(fmt.Sprintf("Migration index %d registered multiple times", index))
	}
	m.migrations[index] = ops
}

func (m *MigrationStore) applyAll(db *schema.Schema) {
	var i int64 = 1
	for {
		ops, ok := m.migrations[i]
		if !ok {
			break
		}

		for _, op := range ops {
			op.Apply(db)
		}

		i++
	}
}

func (m *MigrationStore) nextFree() int64 {
	var i int64 = 1
	for {
		_, ok := m.migrations[i]
		if !ok {
			break
		}

		i++
	}
	return i
}
