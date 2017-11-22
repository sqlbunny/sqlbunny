{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $colDefs := sqlColDefinitions .Table.Columns .Table.PKey.Columns -}}
{{- $pkNames := $colDefs.Names | stringMap .StringFuncs.camelCase | stringMap .StringFuncs.replaceReserved -}}
{{- $pkArgs := joinSlices " " $pkNames $colDefs.Types | join ", " -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
// {{$tableNameSingular}}Exists checks if the {{$tableNameSingular}} row exists.
func {{$tableNameSingular}}Exists(ctx context.Context, {{$pkArgs}}) (bool, error) {
	var exists bool
	{{if eq .DriverName "mssql" -}}
	sql := "select case when exists(select top(1) 1 from {{$schemaTable}} where {{if .Dialect.IndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}}) then 1 else 0 end"
	{{- else -}}
	sql := "select exists(select 1 from {{$schemaTable}} where {{if .Dialect.IndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}} limit 1)"
	{{- end}}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, {{$pkNames | join ", "}})
	}

	row := boil.QueryRowContext(ctx, sql, {{$pkNames | join ", "}})

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "{{.PkgName}}: unable to check if {{.Table.Name}} exists")
	}

	return exists, nil
}
