
import (
    "bytes"
    "database/sql/driver"
    "encoding/json"

    "github.com/sqlbunny/sqlbunny/runtime/bunny"
    "github.com/sqlbunny/sqlbunny/types/null/convert"
)

{{- $dot := . -}}
{{- $modelName := .Struct.Name | titleCase -}}
{{- $modelNameCamel := .Struct.Name | camelCase -}}

// {{$modelName}} is an object representing the database model.
type {{$modelName}} struct {
	{{range $field := .Struct.Fields }}
	{{titleCase $field.Name}} {{goType $field.GoType}} `{{$field.GenerateTags}}`
	{{- end -}}
}
