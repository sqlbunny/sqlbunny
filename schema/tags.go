package schema

import (
	"bytes"
	"errors"
)

var (
	errTagSyntax      = errors.New("bad syntax for struct tag pair")
	errTagKeySyntax   = errors.New("bad syntax for struct tag key")
	errTagValueSyntax = errors.New("bad syntax for struct tag value")

	errKeyNotSet      = errors.New("tag key does not exist")
	errTagNotExist    = errors.New("tag does not exist")
	errTagKeyMismatch = errors.New("mismatch between key and tag.key")
)

// Tags represent a set of tags from a single struct field
type Tags map[string]string

func (t Tags) String() string {
	var buf bytes.Buffer
	for key, value := range t {
		buf.WriteString(key)
		buf.WriteString(":\"")
		buf.WriteString(value)
		buf.WriteString("\" ")
	}
	return buf.String()
}
