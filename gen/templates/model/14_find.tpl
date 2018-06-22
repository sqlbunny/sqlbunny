{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $colDefs := sqlColDefinitions .Model.Columns .Model.PrimaryKey.Columns -}}
{{- $pkNames := $colDefs.Names | stringMap .StringFuncs.camelCase | stringMap .StringFuncs.replaceReserved -}}
{{- $pkArgs := joinSlices " " $pkNames $colDefs.Types | join ", "}}
// Find{{$modelNameSingular}} retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all fields.
func Find{{$modelNameSingular}}(ctx context.Context, {{$pkArgs}}, selectCols ...string) (*{{$modelNameSingular}}, error) {
	{{$varNameSingular}}Obj := &{{$modelNameSingular}}{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"SELECT %s FROM {{.Model.Name | .SchemaModel}} WHERE {{if .Dialect.IndexPlaceholders}}{{whereClause .LQ .RQ 1 .Model.PrimaryKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Model.PrimaryKey.Columns}}{{end}}", sel,
	)

	q := queries.Raw(ctx, query, {{$pkNames | join ", "}})

	err := q.Bind({{$varNameSingular}}Obj)
	if err != nil {
		return nil, errors.Wrap(err, "{{.PkgName}}: unable to select from {{.Model.Name}}")
	}

	return {{$varNameSingular}}Obj, nil
}
