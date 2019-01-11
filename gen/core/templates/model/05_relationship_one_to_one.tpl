{{- if .Model.IsJoinModel -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Model.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $dot.Models $dot.Model . -}}
		{{- $varNameSingular := .ForeignModel | singular | camelCase}}
// {{$txt.Function.NameGo}}G pointed to by the foreign key.

// {{$txt.Function.NameGo}} pointed to by the foreign key.
func (o *{{$txt.LocalModel.NameGo}}) {{$txt.Function.NameGo}}(mods ...qm.QueryMod) ({{$varNameSingular}}Query) {
	queryMods := []qm.QueryMod{
		qm.Where("{{$txt.ForeignModel.ColumnName}}=?", o.{{$txt.LocalModel.ColumnNameGo}}),
	}

	queryMods = append(queryMods, mods...)

	query := {{$txt.ForeignModel.NamePluralGo}}(queryMods...)
	queries.SetFrom(query.Query, "{{.ForeignModel | schemaModel}}")

	return query
}
{{- end -}}
{{- end -}}
