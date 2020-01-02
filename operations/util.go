package operations

import (
	"bytes"
	"fmt"
)

func esc(s string) string {
	return fmt.Sprintf("%#v", s)
}

func dumpBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
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
