{{- define "relationship_to_one_struct_helper" -}}
{{- end -}}

{{- $dot := . -}}
{{- $tableNameSingular := .Table.Name | singular -}}
{{- $modelName := $tableNameSingular | titleCase -}}
{{- $modelNameCamel := $tableNameSingular | camelCase -}}

// {{$modelNameCamel}}JSON is an object representing the JSON serialized form of {{$modelName}}.
type {{$modelNameCamel}}JSON struct {
	{{range $column := .Table.Columns }}
	{{- if eq $dot.StructTagCasing "camel"}}
	{{titleCase $column.Name}} {{$column.Type}} `{{generateTags $dot.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name | camelCase}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name | camelCase}}" yaml:"{{$column.Name | camelCase}}{{if $column.Nullable}},omitempty{{end}}"`
	{{- else -}}
	{{titleCase $column.Name}} {{$column.Type}} `{{generateTags $dot.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name}}" yaml:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}"`
	{{end -}}
	{{end -}}
	{{range .Table.FKeys -}}
	{{- $txt := txtsFromFKey $dot.Tables $dot.Table . -}}
	{{$txt.Function.NameGo}} *{{.ForeignTable | camelCase}}JSON `json:"{{$txt.Function.Name}},omitempty"`
	{{end -}}
}

func (o *{{$modelName}}) JSON() *{{$modelNameCamel}}JSON {
    if o == nil {
        return nil
    }

    res := &{{$modelNameCamel}}JSON{
        {{range $column := .Table.Columns -}}
        {{titleCase $column.Name}}: o.{{titleCase $column.Name}},
        {{end -}}
    }

    if o.R != nil {
        {{range .Table.FKeys -}}
        {{- $txt := txtsFromFKey $dot.Tables $dot.Table . -}}
        res.{{$txt.Function.NameGo}} = o.R.{{$txt.Function.NameGo}}.JSON()
        {{end -}}
    }
    return res
}

func (o {{$modelName}}Slice) JSON() []*{{$modelNameCamel}}JSON {
    if o == nil {
        return nil
    }
    res := make([]*{{$modelNameCamel}}JSON, len(o))
    for i := range o {
        res[i] = o[i].JSON()
    }
    return res
}
