package schema

import (
	"bytes"
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
