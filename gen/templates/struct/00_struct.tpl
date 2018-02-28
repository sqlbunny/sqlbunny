{{- $dot := . -}}
{{- $modelName := .Struct.Name | titleCase -}}
{{- $modelNameCamel := .Struct.Name | camelCase -}}

// {{$modelName}} is an object representing the database model.
type {{$modelName}} struct {
	{{range $field := .Struct.Fields }}
	{{- if eq $dot.StructTagCasing "camel"}}
	{{titleCase $field.Name}} {{$field.TypeGo}} `{{generateTags $dot.Tags $field.Name}}boil:"{{$field.Name}}" json:"{{$field.Name | camelCase}}{{if $field.Nullable}},omitempty{{end}}" toml:"{{$field.Name | camelCase}}" yaml:"{{$field.Name | camelCase}}{{if $field.Nullable}},omitempty{{end}}"`
	{{- else -}}
	{{titleCase $field.Name}} {{$field.TypeGo}} `{{generateTags $dot.Tags $field.Name}}boil:"{{$field.Name}}" json:"{{$field.Name}}{{if $field.Nullable}},omitempty{{end}}" toml:"{{$field.Name}}" yaml:"{{$field.Name}}{{if $field.Nullable}},omitempty{{end}}"`
	{{end -}}
	{{end -}}
}
