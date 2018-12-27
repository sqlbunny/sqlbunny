{{- if .Model.IsJoinModel -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Model.ForeignKeys -}}
		{{- $txt := txtsFromFKey $dot.Models $dot.Model . -}}
		{{- $varNameSingular := .ForeignModel | singular | camelCase}}

// {{$txt.Function.NameGo}} pointed to by the foreign key.
func (o *{{$txt.LocalModel.NameGo}}) {{$txt.Function.NameGo}}(ctx context.Context, mods ...qm.QueryMod) ({{$varNameSingular}}Query) {
	queryMods := []qm.QueryMod{
		qm.Where("{{$txt.ForeignModel.ColumnName}}=?", o.{{$txt.LocalModel.ColumnNameGo}}),
	}

	queryMods = append(queryMods, mods...)

	query := {{$txt.ForeignModel.NamePluralGo}}(ctx, queryMods...)
	queries.SetFrom(query.Query, "{{.ForeignModel | schemaModel}}")

	return query
}
{{- end -}}
{{- end -}}
