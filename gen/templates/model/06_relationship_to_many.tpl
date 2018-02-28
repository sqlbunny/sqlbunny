{{- if .Model.IsJoinModel -}}
{{- else -}}
	{{- $dot := . -}}
	{{- $model := .Model -}}
	{{- range .Model.ToManyRelationships -}}
		{{- $varNameSingular := .ForeignModel | singular | camelCase -}}
		{{- $txt := txtsFromToMany $dot.Models $model . -}}
		{{- $schemaForeignModel := .ForeignModel | $dot.SchemaModel}}

// {{$txt.Function.NameGo}} retrieves all the {{.ForeignModel | singular}}'s {{$txt.ForeignModel.NameHumanReadable}} with an executor
{{- if not (eq $txt.Function.NameGo $txt.ForeignModel.NamePluralGo)}} via {{.ForeignColumn}} field{{- end}}.
func (o *{{$txt.LocalModel.NameGo}}) {{$txt.Function.NameGo}}(ctx context.Context, mods ...qm.QueryMod) {{$varNameSingular}}Query {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

		{{if .ToJoinModel -}}
	queryMods = append(queryMods,
		{{$schemaJoinModel := .JoinModel | $.SchemaModel -}}
		qm.InnerJoin("{{$schemaJoinModel}} on {{$schemaForeignModel}}.{{.ForeignColumn | $dot.Quotes}} = {{$schemaJoinModel}}.{{.JoinForeignColumn | $dot.Quotes}}"),
		qm.Where("{{$schemaJoinModel}}.{{.JoinLocalField | $dot.Quotes}}=?", o.{{$txt.LocalModel.ColumnNameGo}}),
	)
		{{else -}}
	queryMods = append(queryMods,
		qm.Where("{{$schemaForeignModel}}.{{.ForeignColumn | $dot.Quotes}}=?", o.{{$txt.LocalModel.ColumnNameGo}}),
	)
		{{end}}

	query := {{$txt.ForeignModel.NamePluralGo}}(ctx, queryMods...)
	queries.SetFrom(query.Query, "{{$schemaForeignModel}}")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"{{$schemaForeignModel}}.*"})
	}

	return query
}

{{end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* if isJoinModel */ -}}
