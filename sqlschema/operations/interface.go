package operations

import (
	"github.com/sqlbunny/sqlbunny/sqlschema/schema"
)

type Operation interface {
	// GetSQL returns the SQL for the operation. transaction reports whether the
	// migration is running inside a transaction.
	GetSQL(transaction bool) string
	Apply(d *schema.Database) error
}
