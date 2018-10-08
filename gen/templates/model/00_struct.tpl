{{ import "bytes" "bytes" }}
{{ import "context" "context" }}
{{ import "sql" "database/sql" }}
{{ import "json" "encoding/json" }}
{{ import "fmt" "fmt" }}
{{ import "reflect" "reflect" }}
{{ import "strings" "strings" }}
{{ import "sync" "sync" }}
{{ import "time" "time" }}
{{ import "bunny" "github.com/KernelPay/sqlbunny/bunny" }}
{{ import "queries" "github.com/KernelPay/sqlbunny/bunny/queries" }}
{{ import "qm" "github.com/KernelPay/sqlbunny/bunny/qm" }}
{{ import "strmangle" "github.com/KernelPay/sqlbunny/bunny/strmangle" }}
{{ import "errors" "github.com/pkg/errors" }}

{{- $dot := . -}}
{{- $modelNameSingular := .Model.Name | singular -}}
{{- $modelName := $modelNameSingular | titleCase -}}
{{- $modelNameCamel := $modelNameSingular | camelCase -}}

// {{$modelName}} is an object representing the database model.
type {{$modelName}} struct {
	{{range $field := .Model.Fields }}
    {{titleCase $field.Name}} {{typeGo $field.TypeGo}} `{{$field.GenerateTags}}`
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
	{{range .Model.ForeignKeys -}}
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
