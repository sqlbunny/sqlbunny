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
{{ import "errors" "github.com/pkg/errors" }}

{{- $dot := . -}}
{{- $modelNameSingular := .Model.Name | singular -}}
{{- $modelName := $modelNameSingular | titleCase -}}
{{- $modelNameCamel := $modelNameSingular | camelCase -}}

// {{$modelName}} is an object representing the database model.
type {{$modelName}} struct {
	{{range $field := .Model.Fields }}
    {{titleCase $field.Name}} {{goType $field.GoType}} `{{$field.GenerateTags}}`
	{{- end -}}
	{{- if .Model.IsJoinModel -}}
	{{- else}}
	R *{{$modelNameCamel}}R `json:"-" toml:"-" yaml:"-"`
	L {{$modelNameCamel}}L `json:"-" toml:"-" yaml:"-"`
	{{end -}}
}

var {{$modelName}}Columns = struct {
	{{range $column := .Model.Columns -}}
	{{titleCase $column.Name}} string
	{{end -}}
}{
	{{range $column := .Model.Columns -}}
	{{titleCase $column.Name | }}: "{{$column.Name}}",
	{{end -}}
}

{{- if .Model.IsJoinModel -}}
{{- else}}
// {{$modelNameCamel}}R is where relationships are stored.
type {{$modelNameCamel}}R struct {
	{{range .Model.SingleColumnForeignKeys -}}
	{{- $txt := txtsFromFKey $dot.Models $dot.Model . -}}
	{{$txt.Function.NameGo}} *{{$txt.ForeignModel.NameGo}}
	{{end -}}

	{{range .Model.ToOneRelationships -}}
	{{- $txt := txtsFromOneToOne $dot.Models $dot.Model . -}}
	{{$txt.Function.NameGo}} *{{$txt.ForeignModel.NameGo}}
	{{end -}}

	{{range .Model.ToManyRelationships -}}
	{{- $txt := txtsFromToMany $dot.Models $dot.Model . -}}
	{{$txt.Function.NameGo}} {{$txt.ForeignModel.Slice}}
	{{end -}}{{/* range tomany */}}
}

// {{$modelNameCamel}}L is where Load methods for each relationship are stored.
type {{$modelNameCamel}}L struct{}
{{end -}}
