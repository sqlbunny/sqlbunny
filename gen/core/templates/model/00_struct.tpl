{{ import "bytes" "bytes" }}
{{ import "context" "context" }}
{{ import "sql" "database/sql" }}
{{ import "json" "encoding/json" }}
{{ import "fmt" "fmt" }}
{{ import "reflect" "reflect" }}
{{ import "sort" "sort" }}
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
	R {{$modelNameCamel}}R `json:"-" toml:"-" yaml:"-"`
	L {{$modelNameCamel}}L `json:"-" toml:"-" yaml:"-"`
}

type {{$modelName}}Column string

var {{$modelName}}Columns = struct {
	{{range $name, $column := .Model.Table.Columns -}}
	{{titleCase $name}} {{$modelName}}Column
	{{end -}}
}{
	{{range $name, $column := .Model.Table.Columns -}}
	{{titleCase $name | }}: "{{$name}}",
	{{end -}}
}

// {{$modelNameCamel}}R is where relationships are stored.
// Fields are private; access them via the generated getter methods, which
// panic if the relationship hasn't been loaded. Use qm.Load("rel_name")
// in your query to load.
type {{$modelNameCamel}}R struct {
	{{range .Model.Relationships -}}
	{{ .Name | camelCase }}Loaded bool
	{{- if .ToMany }}
	{{ .Name | camelCase }} {{ .ForeignModel | titleCase}}Slice
	{{- else }}
	{{ .Name | camelCase }} *{{ .ForeignModel | titleCase}}
	{{- end }}
	{{end -}}
}

{{range .Model.Relationships}}
// {{ .Name | titleCase }} returns the loaded {{ .Name }} relationship.
// Panics if the relationship was not loaded via qm.Load("{{ .Name }}").
func (r *{{$modelNameCamel}}R) {{ .Name | titleCase }}() {{if .ToMany}}{{ .ForeignModel | titleCase}}Slice{{else}}*{{ .ForeignModel | titleCase}}{{end}} {
	if !r.{{ .Name | camelCase }}Loaded {
		panic(`sqlbunny: relationship "{{ .Name | titleCase }}" on {{$modelName}} not loaded; use qm.Load("{{ .Name }}")`)
	}
	return r.{{ .Name | camelCase }}
}
{{end}}

// {{$modelNameCamel}}L is where Load methods for each relationship are stored.
type {{$modelNameCamel}}L struct{}
