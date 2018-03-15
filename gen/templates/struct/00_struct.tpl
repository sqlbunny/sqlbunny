{{- $dot := . -}}
{{- $modelName := .Struct.Name | titleCase -}}
{{- $modelNameCamel := .Struct.Name | camelCase -}}

// {{$modelName}} is an object representing the database model.
type {{$modelName}} struct {
	{{range $field := .Struct.Fields }}
	{{titleCase $field.Name}} {{$field.TypeGo}} `{{$field.GenerateTags}}`
	{{end -}}
}
