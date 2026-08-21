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
		if err := appendArrayElementOpt(&b, e, true); err != nil {
			return nil, err
		}
	}
	b.WriteByte('}')
	return b.String(), nil
}

func appendArrayElementOpt(b *strings.Builder, e any, bytesAsHex bool) error {
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
		if !bytesAsHex {
			appendArrayQuoted(b, string(v))
			break
		}
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

// ColumnArray encodes the values of one column as a single array parameter
// for "unnest(?::T[])", where T is the returned paramType.
//
// sqlType is the column's Postgres type. For array-typed columns the values
// (array literals, as returned by their Valuer) are sent as text[] and the
// caller casts each element back to the real type. []byte values are hex
// encoded only for bytea columns; for other columns (json, text) they are
// sent as text.
func ColumnArray(sqlType string, vals []any) (arr driver.Valuer, paramType string) {
	if strings.HasSuffix(sqlType, "[]") {
		return anyArray{vals: vals, bytesAsHex: false}, "text[]"
	}
	return anyArray{vals: vals, bytesAsHex: sqlType == "bytea"}, sqlType + "[]"
}

type anyArray struct {
	vals       []any
	bytesAsHex bool
}

func (a anyArray) Value() (driver.Value, error) {
	var b strings.Builder
	b.WriteByte('{')
	for i, e := range a.vals {
		if i > 0 {
			b.WriteByte(',')
		}
		if err := appendArrayElementOpt(&b, e, a.bytesAsHex); err != nil {
			return nil, err
		}
	}
	b.WriteByte('}')
	return b.String(), nil
}

// BulkValues builds the source of a multi-row INSERT with one array parameter
// per column, so the parameter count is len(colTypes), not rows*columns:
//
//	SELECT u.c0, u.c1::bytea[] FROM unnest($1::bytea[], $2::text[]) AS u(c0, c1)
//
// colTypes are the Postgres types of the columns, rows the values in the same
// column order. Row order is preserved.
func BulkValues(colTypes []string, rows [][]any) (string, []any) {
	n := len(colTypes)
	colVals := make([][]any, n)
	for i := range colVals {
		colVals[i] = make([]any, len(rows))
	}
	for r, row := range rows {
		for c := 0; c < n; c++ {
			colVals[c][r] = row[c]
		}
	}

	args := make([]any, n)
	sel := make([]string, n)
	params := make([]string, n)
	names := make([]string, n)
	for c, t := range colTypes {
		arr, paramType := ColumnArray(t, colVals[c])
		args[c] = arr
		names[c] = fmt.Sprintf("c%d", c)
		sel[c] = "u." + names[c]
		if strings.HasSuffix(t, "[]") {
			sel[c] += "::" + t
		}
		params[c] = placeholder(c+1) + "::" + paramType
	}
	return fmt.Sprintf("SELECT %s FROM unnest(%s) AS u(%s)",
		strings.Join(sel, ", "), strings.Join(params, ", "), strings.Join(names, ", ")), args
}
