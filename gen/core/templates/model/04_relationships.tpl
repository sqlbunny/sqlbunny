{{- $dot := . -}}
{{- $model := .Model -}}
{{- $modelName := .Model.Name | titleCase -}}
{{- $modelNameCamel := .Model.Name | camelCase -}}
{{- $modelNamePlural := .Model.Name | plural | titleCase -}}

{{ range .Model.Relationships -}}

{{- $relationship := . }}
{{- $relationshipName := .Name | titleCase}}
{{- $relationshipNameCamel := .Name | camelCase}}


{{ $foreignModel := index $dot.Schema.Models .ForeignModel }}
{{- $foreignModelName := .ForeignModel | titleCase}}
{{- $foreignModelNameCamel := .ForeignModel | camelCase}}
{{- $foreignModelNamePlural := .ForeignModel | plural | titleCase -}}

func (o *{{$modelName}}) {{$relationshipName}}(mods ...qm.QueryMod) ({{$foreignModelNameCamel}}Query) {
	queryMods := []qm.QueryMod{
		{{if .IsJoinModel -}}
		qm.InnerJoin("{{.JoinModel | schemaModel }} ON {{joinOnClause .JoinModel .JoinForeignFields .ForeignModel .ForeignFields}}"),
		qm.Where("{{joinWhereClause 0 .JoinModel .JoinLocalFields}}" {{range .LocalFields}}, o.{{. | titleCasePath}}{{end}}),
		{{if .JoinWhere -}}
		{{- $schemaModel := .JoinModel | schemaModel }}
		qm.Where("{{replaceAll .JoinWhere "$join" $schemaModel}}"),
		{{- end }}
		{{ else if .IsArray }}
		{{- $localArray := index .LocalFields 0 -}}
		{{- $foreignCol := index .ForeignFields 0 }}
		qm.Where("\"{{$foreignCol.SQLName}}\" = ANY(?)", o.{{$localArray | titleCasePath}}),
		{{ else }}
		qm.Where("{{whereClause 0 .ForeignFields}}" {{range .LocalFields}}, o.{{. | titleCasePath}}{{end}}),
		{{- end }}
		{{if .ForeignWhere -}}
		{{- $schemaModel := .ForeignModel | schemaModel }}
		qm.Where("{{replaceAll .ForeignWhere "$foreign" $schemaModel}}"),
		{{- end }}
		{{if .ForeignOrderBy -}}
		qm.OrderBy("{{.ForeignOrderBy}}"),
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
func ({{$modelNameCamel}}L) Load{{$relationshipName}}(ctx context.Context, slice []*{{$modelName}}, mods ...qm.QueryMod) error {
	// Mark loaded and reset any previous value so re-loads start from a clean slate.
	for _, obj := range slice {
		obj.R.{{$relationshipNameCamel}}Loaded = true
		obj.R.{{$relationshipNameCamel}} = nil
	}

	{{if .IsArray}}
	{{- $localArray := index .LocalFields 0 -}}
	{{- $foreignCol := index .ForeignFields 0 -}}
	{{- $fcol := $foreignModel.FindField $foreignCol }}

	// Collect unique foreign IDs across all local arrays.
	seen := make(map[{{ goType $fcol.GoType }}]struct{})
	var allIDs []{{ goType $fcol.GoType }}
	for _, obj := range slice {
		for _, id := range obj.{{$localArray | titleCasePath}} {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			allIDs = append(allIDs, id)
		}
	}
	if len(allIDs) == 0 {
		return nil
	}

	query := NewQuery(
		qm.Select("f.*"),
		qm.From("{{.ForeignModel | schemaModel}} AS f"),
		qm.Where("f.\"{{$foreignCol.SQLName}}\" = ANY(?::{{sqlType $fcol}}[])", queries.Array(allIDs)),
		{{if .ForeignWhere -}}
		qm.Where("{{replaceAll .ForeignWhere "$foreign" "f"}}"),
		{{- end }}
		{{if .ForeignOrderBy -}}
		qm.OrderBy("{{.ForeignOrderBy}}"),
		{{- end }}
	)
	qm.Apply(query, mods...)

	var resultSlice []*{{$foreignModelName}}
	if err := query.Bind(ctx, &resultSlice); err != nil {
		return errors.Errorf("failed to bind eager loaded slice {{$foreignModelName}}: %w", err)
	}

	{{ hook $dot "after_select_slice_noreturn" "resultSlice" $foreignModel }}

	if len(resultSlice) == 0 {
		return nil
	}

	// Index foreign rows by their key for O(1) lookup, preserving order.
	foreignByKey := make(map[{{ goType $fcol.GoType }}]*{{$foreignModelName}}, len(resultSlice))
	for _, foreign := range resultSlice {
		foreignByKey[foreign.{{$foreignCol | titleCasePath}}] = foreign
	}

	for _, local := range slice {
		for _, id := range local.{{$localArray | titleCasePath}} {
			if foreign, ok := foreignByKey[id]; ok {
				local.R.{{$relationshipNameCamel}} = append(local.R.{{$relationshipNameCamel}}, foreign)
			}
		}
	}
	{{else}}
	if len(slice) == 0 {
		return nil
	}

	// Collect unique keys. One typed array per key column, sent as one
	// parameter each, so the parameter count does not grow with the slice.
	{{- range $i, $c := .LocalFields }}
	{{- $lcol := $model.FindField $c }}
	keys{{$i}} := make([]{{ goType $lcol.GoType }}, 0, len(slice))
	{{- end }}
	{{- if eq (len .LocalFields) 1 }}
	{{- $lcol := $model.FindField (index .LocalFields 0) }}
	seen := make(map[{{ goType $lcol.GoType }}]struct{}, len(slice))
	for _, obj := range slice {
		k := obj.{{ index .LocalFields 0 | titleCasePath }}
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		keys0 = append(keys0, k)
	}
	{{- else }}
	type keyType struct {
		{{- range $i, $c := .LocalFields }}
		{{- $lcol := $model.FindField $c }}
		K{{$i}} {{ goType $lcol.GoType }}
		{{- end }}
	}
	seen := make(map[keyType]struct{}, len(slice))
	for _, obj := range slice {
		k := keyType{ {{- range $i, $c := .LocalFields }}{{if $i}}, {{end}}obj.{{$c | titleCasePath}}{{ end -}} }
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		{{- range $i, $c := .LocalFields }}
		keys{{$i}} = append(keys{{$i}}, k.K{{$i}})
		{{- end }}
	}
	{{- end }}

	{{if .IsJoinModel }}
	{{ $joinModel := index $dot.Schema.Models .JoinModel }}
	{{- $joinModelName := .JoinModel | titleCase}}
	{{- $joinModelNameCamel := .JoinModel | camelCase}}

	{{- if eq (len .LocalFields) 1 }}
	{{- $jcol := $joinModel.FindField (index .JoinLocalFields 0) }}
	where := "{{ whereInClause "j" .JoinLocalFields }} = ANY(?::{{sqlType $jcol}}[])"
	{{- else }}
	where := "{{ whereInClause "j" .JoinLocalFields }} IN (SELECT * FROM unnest(
		{{- range $i, $c := .JoinLocalFields }}{{if $i}}, {{end}}?::{{ sqlType ($joinModel.FindField $c) }}[]{{ end -}}
	))"
	{{- end }}
	query := NewQuery(
		qm.Select(
			{{ range $i, $c := $foreignModel.Table.Columns -}}"f.{{$i}}",{{end}}
			{{ range $i, $c := .JoinLocalFields -}}{{if $i}},{{end}} "j.{{$c.SQLName}}"{{end}},
		),
		qm.From("{{.ForeignModel | schemaModel}} AS f"),
		qm.InnerJoin("{{.JoinModel | schemaModel }} AS j ON {{joinOnClause "j" .JoinForeignFields "f" .ForeignFields}}"),
		qm.Where(where{{ range $i, $c := .LocalFields }}, queries.Array(keys{{$i}}){{ end }}),
		{{if .JoinWhere -}}
		qm.Where("{{replaceAll .JoinWhere "$join" "j"}}"),
		{{- end }}
		{{if .ForeignWhere -}}
		qm.Where("{{replaceAll .ForeignWhere "$foreign" "f"}}"),
		{{- end }}
		{{if .ForeignOrderBy -}}
		qm.OrderBy("{{.ForeignOrderBy}}"),
		{{- end }}
	)
	qm.Apply(query, mods...)
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
				local.R.{{$relationshipNameCamel}} = append(local.R.{{$relationshipNameCamel}}, &joined.F)
				{{else}}
				local.R.{{$relationshipNameCamel}} = &joined.F
				{{end}}
			}
		}
	}
	{{else}}
	{{- if eq (len .LocalFields) 1 }}
	{{- $fcol := $foreignModel.FindField (index .ForeignFields 0) }}
	where := "{{ whereInClause "f" .ForeignFields }} = ANY(?::{{sqlType $fcol}}[])"
	{{- else }}
	where := "{{ whereInClause "f" .ForeignFields }} IN (SELECT * FROM unnest(
		{{- range $i, $c := .ForeignFields }}{{if $i}}, {{end}}?::{{ sqlType ($foreignModel.FindField $c) }}[]{{ end -}}
	))"
	{{- end }}
	query := NewQuery(
		qm.Select("f.*"),
		qm.From("{{.ForeignModel | schemaModel}} AS f"),
		qm.Where(where{{ range $i, $c := .LocalFields }}, queries.Array(keys{{$i}}){{ end }}),
		{{if .ForeignWhere -}}
		qm.Where("{{replaceAll .ForeignWhere "$foreign" "f"}}"),
		{{- end }}
		{{if .ForeignOrderBy -}}
		qm.OrderBy("{{.ForeignOrderBy}}"),
		{{- end }}
	)
	qm.Apply(query, mods...)

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
				local.R.{{$relationshipNameCamel}} = append(local.R.{{$relationshipNameCamel}}, foreign)
				{{else}}
				local.R.{{$relationshipNameCamel}} = foreign
				break
				{{end}}
			}
		}
	}
	{{end}}
	{{end}}

	return nil
}

{{ end -}}
