{{- $dot := . -}}
{{- $enumName := .Enum.Name | titleCase -}}
{{- $enumNameCamel := .Enum.Name | camelCase -}}

import (
    "bytes"
    "database/sql/driver"
    "encoding/json"
    "bytes"

    "github.com/kernelpayments/sqlbunny/bunny"
    "github.com/kernelpayments/sqlbunny/types/null/convert"
)

// {{$enumName}} is an enum type.
type {{$enumName}} {{typeGo .Enum.BaseTypeGo}}

const (
    {{- range $index, $choice := .Enum.Choices }}
    {{$enumName}}{{$choice | titleCase}} = {{$enumName}}({{$index}})
    {{- end}}
)

var {{$enumNameCamel}}Values = map[string]{{$enumName}}{
    {{- range $index, $choice := .Enum.Choices }}
    "{{$choice}}": {{$enumName}}({{$index}}),
    {{- end}}
}

var {{$enumNameCamel}}Names = map[{{$enumName}}]string{
    {{- range $index, $choice := .Enum.Choices }}
    {{$enumName}}({{$index}}): "{{$choice}}",
    {{- end}}
}

func (o {{$enumName}}) String() string {
    return {{$enumNameCamel}}Names[o]
}

func {{$enumName}}FromString(s string) ({{$enumName}}, error) {
    var o {{$enumName}}
    err := o.UnmarshalText([]byte(s))
    return o, err
}

// MarshalText implements encoding/text TextMarshaler interface.
func (o {{$enumName}}) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

// UnmarshalText implements encoding/text TextUnmarshaler interface.
func (o *{{$enumName}}) UnmarshalText(text []byte) error {
	val, ok := {{$enumNameCamel}}Values[string(text)]
	if !ok {
        return &bunny.InvalidEnumError{Value: text, Type: "{{$enumName}}"}
	}
	*o = val
	return nil
}
