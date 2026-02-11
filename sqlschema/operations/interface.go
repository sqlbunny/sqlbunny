package operations

import (
	"github.com/sqlbunny/sqlbunny/sqlschema/schema"
)

type Operation interface {
	GetSQL() string
	Apply(d *schema.Database) error
}
