package operations

import (
	"github.com/sqlbunny/sqlschema/schema"
)

type Operation interface {
	GetSQL() string
	Apply(d *schema.Database) error
}
