package migration

import (
	"github.com/sqlbunny/sqlbunny/runtime/migration"
	"github.com/sqlbunny/sqlschema/schema"
)

func ApplyMigration(m *migration.Migration, s *schema.Schema) error {
	for _, o := range m.Operations {
		if err := o.Apply(s); err != nil {
			return err
		}
	}
	return nil
}
