{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *{{$tableNameSingular}}) Reload(ctx context.Context) error {
	ret, err := Find{{$tableNameSingular}}(ctx, {{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *{{$tableNameSingular}}Slice) ReloadAll(ctx context.Context) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	{{$varNamePlural}} := {{$tableNameSingular}}Slice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$varNameSingular}}PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT {{$schemaTable}}.* FROM {{$schemaTable}} WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), {{if .Dialect.IndexPlaceholders}}1{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns, len(*o))

	q := queries.Raw(ctx, sql, args...)

	err := q.Bind(&{{$varNamePlural}})
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to reload all in {{$tableNameSingular}}Slice")
	}

	*o = {{$varNamePlural}}

	return nil
}
