{{- define "relationship_to_one_struct_helper" -}}
{{- end -}}

{{- $dot := . -}}
{{- $modelNameSingular := .Model.Name | singular -}}
{{- $modelName := $modelNameSingular | titleCase -}}
{{- $modelNameCamel := $modelNameSingular | camelCase -}}

{{ if .Model.IsStandardModel }}

// {{$modelNameCamel}}JSON is an object representing the JSON serialized form of {{$modelName}}.
type {{$modelNameCamel}}JSON struct {
	{{range $field := .Model.Fields }}
    {{- if not ( $field.HasTag "private" ) }}
	{{ if eq $dot.StructTagCasing "camel" -}}
	{{titleCase $field.Name}} {{$field.TypeGo}} `{{generateTags $dot.Tags $field.Name}}boil:"{{$field.Name}}" json:"{{$field.Name | camelCase}}{{if $field.Nullable}},omitempty{{end}}" toml:"{{$field.Name | camelCase}}" yaml:"{{$field.Name | camelCase}}{{if $field.Nullable}},omitempty{{end}}"`
	{{- else -}}
	{{titleCase $field.Name}} {{$field.TypeGo}} `{{generateTags $dot.Tags $field.Name}}boil:"{{$field.Name}}" json:"{{$field.Name}}{{if $field.Nullable}},omitempty{{end}}" toml:"{{$field.Name}}" yaml:"{{$field.Name}}{{if $field.Nullable}},omitempty{{end}}"`
	{{- end -}}
    {{- end }}
	{{- end }}

	{{range .Model.ForeignKeys -}}
	{{- $txt := txtsFromFKey $dot.Models $dot.Model . -}}
	{{$txt.Function.NameGo}} *{{.ForeignModel | titleCase}} `json:"{{$txt.Function.Name}},omitempty"`
	{{end }}

    CreatedAt time.Time `json:"created_at"`
}


func (o *{{$modelName}}) MarshalJSON() ([]byte, error) {
    if o == nil {
        return []byte("null"), nil
    }

    res := &{{$modelNameCamel}}JSON{
        {{- range $field := .Model.Fields -}}
        {{- if not ( $field.HasTag "private" ) }}
        {{titleCase $field.Name}}: o.{{titleCase $field.Name}},
        {{- end -}}
        {{end }}
    }

    if o.R != nil {
        {{range .Model.ForeignKeys -}}
        {{- $txt := txtsFromFKey $dot.Models $dot.Model . -}}
        res.{{$txt.Function.NameGo}} = o.R.{{$txt.Function.NameGo}}
        {{end -}}
    }

    res.CreatedAt = o.CreatedAt()

    return json.Marshal(res)
}

func (o {{$modelName}}Slice) ToModelSlice() []boil.Model {
    if o == nil {
        return make([]boil.Model, 0)
    }
    res := make([]boil.Model, len(o))
    for i := range o {
        res[i] = o[i]
    }
    return res
}

func (o *{{$modelName}}) CreatedAt() time.Time {
    return o.ID.Time()
}

func (o *{{$modelName}}) GetID() boil.ID {
    return o.ID
}
{{ end }}
