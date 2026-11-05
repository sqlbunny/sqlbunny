package operations

import (
	"fmt"

	"github.com/sqlbunny/sqlbunny/sqlschema/schema"
)

type DropIndex struct {
	SchemaName string
	TableName  string
	IndexName  string
}

// GetSQL uses DROP INDEX CONCURRENTLY outside a transaction. Inside one it
// uses a plain DROP INDEX, since Postgres forbids CONCURRENTLY in a
// transaction block (same as CreateIndex). CONCURRENTLY matters: a plain
// DROP INDEX takes an ACCESS EXCLUSIVE lock on the table, so it waits for
// every running query on the table and blocks all new ones while it waits.
func (o DropIndex) GetSQL(transaction bool) string {
	if transaction {
		return fmt.Sprintf("DROP INDEX %s", sqlName(o.SchemaName, o.IndexName))
	}
	return fmt.Sprintf("DROP INDEX CONCURRENTLY %s", sqlName(o.SchemaName, o.IndexName))
}

func (o DropIndex) Apply(d *schema.Database) error {
	s, ok := d.Schemas[o.SchemaName]
	if !ok {
		return fmt.Errorf("no such schema: %s", o.SchemaName)
	}
	t, ok := s.Tables[o.TableName]
	if !ok {
		return fmt.Errorf("no such table: %s", o.TableName)
	}
	if _, ok := t.Indexes[o.IndexName]; !ok {
		return fmt.Errorf("no such index: %s", o.IndexName)
	}
	delete(t.Indexes, o.IndexName)
	return nil
}
