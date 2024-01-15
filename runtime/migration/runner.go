package migration

import (
	"context"
	"errors"
	"time"

	"github.com/sqlbunny/sqlbunny/runtime/bunny"
)

const (
	checkMigrationsTableSQL  = "SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'migrations'"
	createMigrationsTableSQL = "CREATE TABLE migrations (id text PRIMARY KEY, time timestamptz)"
	insertMigrationSQL       = "INSERT INTO migrations (id, time) VALUES($1, $2)"
	selectMigrationsSQL      = "SELECT id from migrations"
)

func getApplied(ctx context.Context) (map[string]struct{}, error) {

	applied := make(map[string]struct{})
	rows, err := bunny.Query(ctx, selectMigrationsSQL)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied[name] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return applied, nil
}

func (s *Store) ValidateMigrated(ctx context.Context) error {
	applied, err := getApplied(ctx)
	if err != nil {
		return err
	}

	err = s.validateMigrated(applied)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Run(ctx context.Context) error {
	var count int64
	if err := bunny.QueryRow(ctx, checkMigrationsTableSQL).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		if _, err := bunny.Exec(ctx, createMigrationsTableSQL); err != nil {
			return err
		}
	}

	applied, err := getApplied(ctx)
	if err != nil {
		return err
	}

	err = s.validateApplied(applied)
	if err != nil {
		return err
	}

	rows, err := bunny.Query(ctx, selectMigrationsSQL)
	if err != nil {
		return err
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		applied[name] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	heads := s.FindHeads()
	if len(heads) != 1 {
		return errors.New("Found multiple migration heads, you must run 'migration merge' first")
	}
	head := heads[0]

	return s.RunMigration(head, applied, func(m *Migration) error {
		if err := m.Run(ctx); err != nil {
			return err
		}
		if _, err := bunny.Exec(ctx, insertMigrationSQL, m.Name, time.Now()); err != nil {
			return err
		}
		return nil
	})
}
