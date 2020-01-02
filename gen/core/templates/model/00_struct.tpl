{{ import "bytes" "bytes" }}
{{ import "context" "context" }}
{{ import "sql" "database/sql" }}
{{ import "json" "encoding/json" }}
{{ import "fmt" "fmt" }}
{{ import "reflect" "reflect" }}
{{ import "strings" "strings" }}
{{ import "sync" "sync" }}
{{ import "time" "time" }}
{{ import "bunny" "github.com/sqlbunny/sqlbunny/runtime/bunny" }}
{{ import "queries" "github.com/sqlbunny/sqlbunny/runtime/queries" }}
{{ import "qm" "github.com/sqlbunny/sqlbunny/runtime/qm" }}
{{ import "strmangle" "github.com/sqlbunny/sqlbunny/runtime/strmangle" }}
{{ import "errors" "github.com/sqlbunny/errors" }}

{{- $dot := . -}}
{{- $modelName := .Model.Name | titleCase -}}
{{- $modelNameCamel := .Model.Name | camelCase -}}

// {{$modelName}} is an object representing the database model.
type {{$modelName}} struct {
	{{range $field := .Model.Fields }}
    {{titleCase $field.Name}} {{goType $field.GoType}} `{{$field.GenerateTags}}`
	{{- end }}
	R *{{$modelNameCamel}}R `json:"-" toml:"-" yaml:"-"`
	L {{$modelNameCamel}}L `json:"-" toml:"-" yaml:"-"`
}

var {{$modelName}}Columns = struct {
	{{range $name, $column := .Model.Table.Columns -}}
	{{titleCase $name}} string
	{{end -}}
}{
	{{range $name, $column := .Model.Table.Columns -}}
	{{titleCase $name | }}: "{{$name}}",
	{{end -}}
}

// {{$modelNameCamel}}R is where relationships are stored.
type {{$modelNameCamel}}R struct {
	{{range .Model.Relationships -}}
	{{- if .ToMany -}}
	{{ .Name | titleCase }} {{ .ForeignModel | titleCase}}Slice
	{{ else -}}
	{{ .Name | titleCase }} *{{ .ForeignModel | titleCase}}
	{{ end -}}
	{{end -}}
}

// {{$modelNameCamel}}L is where Load methods for each relationship are stored.
type {{$modelNameCamel}}L struct{}
