{{- if .Model.IsJoinModel -}}
{{- else -}}
	{{- $dot := . -}}
	{{- $model := .Model -}}
	{{- range .Model.ToManyRelationships -}}
		{{- $varNameSingular := .ForeignModel | singular | camelCase -}}
		{{- $txt := txtsFromToMany $dot.Models $model . -}}
		{{- $schemaForeignModel := .ForeignModel | schemaModel}}

// {{$txt.Function.NameGo}} retrieves all the {{.ForeignModel | singular}}'s {{$txt.ForeignModel.NameHumanReadable}} with an executor
{{- if not (eq $txt.Function.NameGo $txt.ForeignModel.NamePluralGo)}} via {{.ForeignColumn}} field{{- end}}.
func (o *{{$txt.LocalModel.NameGo}}) {{$txt.Function.NameGo}}(mods ...qm.QueryMod) {{$varNameSingular}}Query {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

		{{if .ToJoinModel -}}
	queryMods = append(queryMods,
		{{$schemaJoinModel := .JoinModel | schemaModel -}}
		qm.InnerJoin("{{$schemaJoinModel}} on {{$schemaForeignModel}}.{{.ForeignColumn | quotes}} = {{$schemaJoinModel}}.{{.JoinForeignColumn | quotes}}"),
		qm.Where("{{$schemaJoinModel}}.{{.JoinLocalColumn | quotes}}=?", o.{{$txt.LocalModel.ColumnNameGo}}),
	)
		{{else -}}
	queryMods = append(queryMods,
		qm.Where("{{$schemaForeignModel}}.{{.ForeignColumn | quotes}}=?", o.{{$txt.LocalModel.ColumnNameGo}}),
	)
		{{end}}

	query := {{$txt.ForeignModel.NamePluralGo}}(queryMods...)
	queries.SetFrom(query.Query, "{{$schemaForeignModel}}")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"{{$schemaForeignModel}}.*"})
	}

	return query
}

{{end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* if isJoinModel */ -}}
