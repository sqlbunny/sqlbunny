{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $dot.Tables $dot.Table . -}}
		{{- $varNameSingular := .ForeignTable | singular | camelCase}}
// {{$txt.Function.Name}}G pointed to by the foreign key.

// {{$txt.Function.Name}} pointed to by the foreign key.
func (o *{{$txt.LocalTable.NameGo}}) {{$txt.Function.Name}}(ctx context.Context, mods ...qm.QueryMod) ({{$varNameSingular}}Query) {
	queryMods := []qm.QueryMod{
		qm.Where("{{$txt.ForeignTable.ColumnName}}=?", o.{{$txt.LocalTable.ColumnNameGo}}),
	}

	queryMods = append(queryMods, mods...)

	query := {{$txt.ForeignTable.NamePluralGo}}(ctx, queryMods...)
	queries.SetFrom(query.Query, "{{.ForeignTable | $dot.SchemaTable}}")

	return query
}
{{- end -}}
{{- end -}}
