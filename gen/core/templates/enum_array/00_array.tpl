{{- $arrayName := .ArrayType.Name | titleCase -}}
{{- $enumName := .ArrayType.Element.Name | titleCase -}}

import (
    "database/sql/driver"

    "github.com/lib/pq"
)

// {{$arrayName}} is a slice of {{$enumName}} backed by a postgres integer[]
// column. It implements sql.Scanner / driver.Valuer by delegating to
// pq.Int32Array, so it interoperates with `= ANY(...)` and `<@ / &&` operators
// natively.
//
// JSON marshaling falls through to the stdlib's handling of []{{$enumName}} —
// {{$enumName}} already implements encoding.TextMarshaler/TextUnmarshaler, so
// encoding/json emits/parses elements as enum-name strings.
type {{$arrayName}} []{{$enumName}}

// Scan implements the sql.Scanner interface.
func (a *{{$arrayName}}) Scan(src any) error {
    var arr pq.Int32Array
    if err := arr.Scan(src); err != nil {
        return err
    }
    if arr == nil {
        *a = nil
        return nil
    }
    out := make({{$arrayName}}, len(arr))
    for i, v := range arr {
        out[i] = {{$enumName}}(v)
    }
    *a = out
    return nil
}

// Value implements the driver.Valuer interface.
func (a {{$arrayName}}) Value() (driver.Value, error) {
    if a == nil {
        return "{}", nil
    }
    arr := make(pq.Int32Array, len(a))
    for i, v := range a {
        arr[i] = int32(v)
    }
    return arr.Value()
}
