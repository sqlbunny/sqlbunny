package operations

import (
	"io"

	"github.com/sqlbunny/sqlschema/schema"
)

type Operation interface {
	GetSQL() string
	Dump(w io.Writer)
	Apply(s *schema.Schema) error
}
