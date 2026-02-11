{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $model := .Model -}}

// Find{{$modelNameSingular}} retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all fields.
func Find{{$modelNameSingular}}(ctx context.Context{{range .Model.PrimaryKey.Fields}}, {{$f := $model.FindField .}}{{$f.Name | camelCase}} {{goType $f.Type.GoType}}{{end}}, selectCols ...{{$modelNameSingular}}Column) (*{{$modelNameSingular}}, error) {
	{{$varNameSingular}}Obj := &{{$modelNameSingular}}{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, columnStrings(selectCols)), ",")
	}
	query := fmt.Sprintf(
		"SELECT %s FROM {{.Model.Name | schemaModel}} WHERE {{if .Dialect.IndexPlaceholders}}{{whereClause .LQ .RQ 1 .Model.PrimaryKey.Fields}}{{else}}{{whereClause .LQ .RQ 0 .Model.PrimaryKey.Fields}}{{end}}", sel,
	)

	q := queries.Raw(query{{range .Model.PrimaryKey.Fields}}, {{$f := $model.FindField .}}{{$f.Name | camelCase}}{{end}})

	err := q.Bind(ctx, {{$varNameSingular}}Obj)
	if err != nil {
		return nil, errors.Errorf("{{.PkgName}}: unable to select from {{.Model.Name}}: %w", err)
	}

	return {{$varNameSingular}}Obj, nil
}
