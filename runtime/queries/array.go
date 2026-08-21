package queries

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Array wraps a slice so it is sent to Postgres as one array parameter.
// Use it with "= ANY(?::type[])" instead of "IN (?, ?, ...)" so the number
// of query parameters does not grow with the slice length.
//
// Elements are converted with their driver.Valuer if they have one, or with
// driver.DefaultParameterConverter otherwise. []byte values are encoded as
// bytea hex literals.
func Array[T any](s []T) driver.Valuer {
	return arrayValuer[T]{s: s}
}

type arrayValuer[T any] struct {
	s []T
}

func (a arrayValuer[T]) Value() (driver.Value, error) {
	var b strings.Builder
	b.WriteByte('{')
	for i, e := range a.s {
		if i > 0 {
			b.WriteByte(',')
		}
		if err := appendArrayElement(&b, e); err != nil {
			return nil, err
		}
	}
	b.WriteByte('}')
	return b.String(), nil
}

func appendArrayElement(b *strings.Builder, e any) error {
	var v driver.Value
	var err error
	if valuer, ok := e.(driver.Valuer); ok {
		v, err = valuer.Value()
	} else {
		v, err = driver.DefaultParameterConverter.ConvertValue(e)
	}
	if err != nil {
		return err
	}

	switch v := v.(type) {
	case nil:
		b.WriteString("NULL")
	case []byte:
		b.WriteString(`"\\x`)
		b.WriteString(hex.EncodeToString(v))
		b.WriteByte('"')
	case string:
		appendArrayQuoted(b, v)
	case bool:
		b.WriteString(strconv.FormatBool(v))
	case int64:
		b.WriteString(strconv.FormatInt(v, 10))
	case float64:
		b.WriteString(strconv.FormatFloat(v, 'g', -1, 64))
	case time.Time:
		appendArrayQuoted(b, v.Format("2006-01-02 15:04:05.999999999Z07:00"))
	default:
		return fmt.Errorf("queries.Array: unsupported element type %T", v)
	}
	return nil
}

func appendArrayQuoted(b *strings.Builder, s string) {
	b.WriteByte('"')
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '"' || c == '\\' {
			b.WriteByte('\\')
		}
		b.WriteByte(c)
	}
	b.WriteByte('"')
}
