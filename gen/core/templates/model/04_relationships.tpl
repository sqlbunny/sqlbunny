{{- $dot := . -}}
{{- $model := .Model -}}
{{- $modelName := .Model.Name | titleCase -}}
{{- $modelNameCamel := .Model.Name | camelCase -}}
{{- $modelNamePlural := .Model.Name | plural | titleCase -}}

{{ range .Model.Relationships -}}

{{- $relationship := . }}
{{- $relationshipName := .Name | titleCase}}


{{ $foreignModel := index $dot.Schema.Models .ForeignModel }}
{{- $foreignModelName := .ForeignModel | titleCase}}
{{- $foreignModelNameCamel := .ForeignModel | camelCase}}
{{- $foreignModelNamePlural := .ForeignModel | plural | titleCase -}}

func (o *{{$modelName}}) {{$relationshipName}}(mods ...qm.QueryMod) ({{$foreignModelNameCamel}}Query) {
	queryMods := []qm.QueryMod{
		{{if .IsJoinModel -}}
		qm.InnerJoin("{{.JoinModel | schemaModel }} ON {{joinOnClause $dot.LQ $dot.RQ .JoinModel .JoinForeignFields .ForeignModel .ForeignFields}}"),
		qm.Where("{{joinWhereClause $dot.LQ $dot.RQ 0 .JoinModel .JoinLocalFields}}" {{range .LocalFields}}, o.{{. | titleCasePath}}{{end}}),
		{{ else }}
		qm.Where("{{whereClause $dot.LQ $dot.RQ 0 .ForeignFields}}" {{range .LocalFields}}, o.{{. | titleCasePath}}{{end}}),
		{{- end }}
	}

	queryMods = append(queryMods, mods...)
	query := {{$foreignModelNamePlural}}(queryMods...)
	queries.SetFrom(query.Query, "{{.ForeignModel | schemaModel}}")
	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"{{.ForeignModel | schemaModel}}.*"})
	}

	return query
}

// Load{{$relationshipName}} allows an eager lookup of values, cached into the
// loaded structs of the objects.
func ({{$modelNameCamel}}L) Load{{$relationshipName}}(ctx context.Context, slice []*{{$modelName}}) error {
	args := make([]interface{}, len(slice)*{{len .LocalFields}})
	for i, obj := range slice {
		if obj.R == nil {
			obj.R = &{{$modelNameCamel}}R{}
		}
		{{ range $i, $c := .LocalFields }}
		args[i*{{len $relationship.LocalFields}} + {{$i}}] = obj.{{$c | titleCasePath}}
		{{ end }}
	}

	{{if .IsJoinModel }}
	{{ $joinModel := index $dot.Schema.Models .JoinModel }}
	{{- $joinModelName := .JoinModel | titleCase}}
	{{- $joinModelNameCamel := .JoinModel | camelCase}}

	where := fmt.Sprintf(
		"{{ whereInClause $dot.LQ $dot.RQ "j" .JoinLocalFields }} in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, len(slice)*{{len .LocalFields}}, 1, {{len .LocalFields}}),
	)
	query := NewQuery(
		qm.Select(
			{{ range $i, $c := $foreignModel.Table.Columns -}}"f.{{$i}}",{{end}}
			{{ range $i, $c := .JoinLocalFields -}}{{if $i}},{{end}} "j.{{$c}}"{{end}},
		),
		qm.From("{{.ForeignModel | schemaModel}} AS f"),
		qm.InnerJoin("{{.JoinModel | schemaModel }} AS j ON {{joinOnClause $dot.LQ $dot.RQ "j" .JoinForeignFields "f" .ForeignFields}}"),
		qm.Where(where, args...),
	)
	type joinStruct struct {
		F {{ $foreignModelName }} `bunny:"f.,bind"`
		J {{ $joinModelName }} `bunny:"j.,bind"`
	}
	var resultSlice []*joinStruct
	if err := query.Bind(ctx, &resultSlice); err != nil {
		return errors.Errorf("failed to bind eager loaded slice {{$foreignModelName}}: %w", err)
	}

	if len(resultSlice) == 0 {
		return nil
	}

	for _, local := range slice {
		for _, joined := range resultSlice {
			if {{ range $i, $lc := .LocalFields -}}
				{{- if $i}} && {{end}}
				{{- $jc := index $relationship.JoinLocalFields $i -}}
				{{- $lcol := $model.FindField $lc -}}
				{{- $jcol := $joinModel.FindField $jc -}}
				{{doCompare (printf "local.%s" ($lc | titleCasePath)) (printf "joined.J.%s" ($jc | titleCasePath)) $lcol $jcol }}
			{{- end }} {
				{{if .ToMany}}
				local.R.{{$relationshipName}} = append(local.R.{{$relationshipName}}, &joined.F)
				{{else}}
				local.R.{{$relationshipName}} = &joined.F
				{{end}}
			}
		}
	}
	{{else}}
	where := fmt.Sprintf(
		"{{ whereInClause $dot.LQ $dot.RQ "f" .ForeignFields }} in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, len(slice)*{{len .LocalFields}}, 1, {{len .LocalFields}}),
	)
	query := NewQuery(
		qm.Select("f.*"),
		qm.From("{{.ForeignModel | schemaModel}} AS f"),
		qm.Where(where, args...),
	)

	var resultSlice []*{{$foreignModelName}}
	if err := query.Bind(ctx, &resultSlice); err != nil {
		return errors.Errorf("failed to bind eager loaded slice {{$foreignModelName}}: %w", err)
	}

	{{ hook $dot "after_select_slice_noreturn" "resultSlice" $foreignModel }}

	if len(resultSlice) == 0 {
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if {{ range $i, $lc := .LocalFields -}}
				{{- if $i}} && {{end}}
				{{- $fc := index $relationship.ForeignFields $i -}}
				{{- $lcol := $model.FindField $lc -}}
				{{- $fcol := $foreignModel.FindField $fc -}}
				{{doCompare (printf "local.%s" ($lc | titleCasePath)) (printf "foreign.%s" ($fc | titleCasePath)) $lcol $fcol }}
			{{- end }} {
				{{if .ToMany}}
				local.R.{{$relationshipName}} = append(local.R.{{$relationshipName}}, foreign)
				{{else}}
				local.R.{{$relationshipName}} = foreign
				break
				{{end}}
			}
		}
	}
	{{end}}

	return nil
}

{{ end -}}
