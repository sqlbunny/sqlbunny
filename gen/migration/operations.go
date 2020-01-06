package migration

import (
	"github.com/sqlbunny/sqlbunny/runtime/migration"
	"github.com/sqlbunny/sqlschema/schema"
)

func ApplyMigration(m *migration.Migration, d *schema.Database) error {
	for _, o := range m.Operations {
		if err := o.Apply(d); err != nil {
			return err
		}
	}
	return nil
}
