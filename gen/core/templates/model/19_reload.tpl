{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $varNamePlural := .Model.Name | plural | camelCase -}}
{{- $schemaModel := .Model.Name | schemaModel}}
// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *{{$modelNameSingular}}) Reload(ctx context.Context) error {
	ret, err := Find{{$modelNameSingular}}(ctx, {{.Model.PrimaryKey.Columns | stringMap .StringFuncs.titleCaseIdentifier | prefixStringSlice "o." | join ", "}})
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key field values
// and overwrites the original object slice with the newly updated slice.
func (o *{{$modelNameSingular}}Slice) ReloadAll(ctx context.Context) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	{{$varNamePlural}} := {{$modelNameSingular}}Slice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$varNameSingular}}PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT {{$schemaModel}}.* FROM {{$schemaModel}} WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), {{if .Dialect.IndexPlaceholders}}1{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns, len(*o))

	q := queries.Raw(ctx, sql, args...)

	err := q.Bind(&{{$varNamePlural}})
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to reload all in {{$modelNameSingular}}Slice")
	}

	*o = {{$varNamePlural}}

	return nil
}
