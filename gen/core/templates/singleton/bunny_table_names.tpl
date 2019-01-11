var ModelNames = struct {
	{{range $model := .Models -}}
	{{titleCase $model.Name}} string
	{{end -}}
}{
	{{range $model := .Models -}}
	{{titleCase $model.Name}}: "{{$model.Name}}",
	{{end -}}
}
