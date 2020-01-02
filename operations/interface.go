package operations

import (
	"bytes"

	"github.com/sqlbunny/sqlschema/schema"
)

type Operation interface {
	GetSQL() string
	Dump(buf *bytes.Buffer)
	Apply(s *schema.Schema) error
}
