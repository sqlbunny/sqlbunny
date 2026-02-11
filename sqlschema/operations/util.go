package operations

import (
	"bytes"
	"fmt"
)

func sqlName(schema, name string) string {
	if schema == "" {
		return fmt.Sprintf("\"%s\"", name)
	}
	return fmt.Sprintf("\"%s\".\"%s\"", schema, name)
}

func columnList(columns []string) string {
	var buf bytes.Buffer
	for i, c := range columns {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("\"")
		buf.WriteString(c)
		buf.WriteString("\"")
	}
	return buf.String()
}
