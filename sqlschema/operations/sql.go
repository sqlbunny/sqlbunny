package operations

import (
	"github.com/sqlbunny/sqlbunny/sqlschema/schema"
)

type SQL struct {
	SQL string
}

func (o SQL) GetSQL() string {
	return o.SQL
}

func (o SQL) Apply(d *schema.Database) error {
	// do nothing.
	return nil
}
