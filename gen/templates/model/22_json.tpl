{{- define "relationship_to_one_struct_helper" -}}
{{- end -}}

{{- $dot := . -}}
{{- $modelNameSingular := .Model.Name | singular -}}
{{- $modelName := $modelNameSingular | titleCase -}}
{{- $modelNameCamel := $modelNameSingular | camelCase -}}

{{ if .Model.IsStandardModel }}

// {{$modelName}}JSON is an object representing the JSON serialized form of {{$modelName}}.
type {{$modelName}}JSON struct {
	{{range $field := .Model.Fields }}
    {{- if not ( $field.HasTag "private" ) }}
	{{titleCase $field.Name}} {{$field.TypeGo}} `json:"{{$field.Name}}{{if $field.Nullable}},omitempty{{end}}" toml:"{{$field.Name}}" yaml:"{{$field.Name}}{{if $field.Nullable}},omitempty{{end}}"`
    {{- end }}
	{{- end }}

	{{range .Model.ForeignKeys -}}
	{{- $txt := txtsFromFKey $dot.Models $dot.Model . -}}
	{{$txt.Function.NameGo}} *{{.ForeignModel | titleCase}} `json:"{{$txt.Function.Name}},omitempty"`
	{{end }}

    CreatedAt time.Time `json:"created_at"`
}


func (o *{{$modelName}}) JSON() *{{$modelName}}JSON {
    if o == nil {
        return nil
    }

    res := &{{$modelName}}JSON{
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

    return res
}

func (o *{{$modelName}}) MarshalJSON() ([]byte, error) {
    return json.Marshal(o.JSON())
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
func (o *{{$modelName}}JSON) GetID() boil.ID {
    return o.ID
}
{{ end }}
