{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $schemaModel := .Model.Name | schemaModel}}
{{- $model := .Model -}}

// {{$modelNameSingular}}Exists checks if the {{$modelNameSingular}} row exists.
func {{$modelNameSingular}}Exists(ctx context.Context{{range .Model.PrimaryKey.Fields}}, {{$f := $model.FindField .}}{{$f.Name | camelCase}} {{goType $f.Type.GoType}}{{end}}, selectCols ...string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from {{$schemaModel}} where {{if .Dialect.IndexPlaceholders}}{{whereClause .LQ .RQ 1 .Model.PrimaryKey.Fields}}{{else}}{{whereClause .LQ .RQ 0 .Model.PrimaryKey.Columns}}{{end}} limit 1)"

	row := bunny.QueryRow(ctx, sql{{range .Model.PrimaryKey.Fields}}, {{$f := $model.FindField .}}{{$f.Name | camelCase}}{{end}})

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Errorf("{{.PkgName}}: unable to check if {{.Model.Name}} exists: %w", err)
	}

	return exists, nil
}
