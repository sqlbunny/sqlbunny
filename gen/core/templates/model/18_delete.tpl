{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $schemaModel := .Model.Name | schemaModel}}
// Delete deletes a single {{$modelNameSingular}} record with an executor.
// Delete will match against the primary key field to find the record to delete.
func (o *{{$modelNameSingular}}) Delete(ctx context.Context) error {
	if o == nil {
	return errors.New("{{.PkgName}}: no {{$modelNameSingular}} provided for delete")
	}

	{{ hook . "before_delete" "o" .Model }}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), {{$varNameSingular}}PrimaryKeyMapping)
	sql := "DELETE FROM {{$schemaModel}} WHERE {{if .Dialect.IndexPlaceholders}}{{whereClause .LQ .RQ 1 .Model.PrimaryKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Model.PrimaryKey.Columns}}{{end}}"

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
	return errors.Wrap(err, "{{.PkgName}}: unable to delete from {{.Model.Name}}")
	}

	{{ hook . "after_delete" "o" .Model }}

	return nil
}

// DeleteAll deletes all matching rows.
func (q {{$varNameSingular}}Query) DeleteAll(ctx context.Context) error {
	if q.Query == nil {
	return errors.New("{{.PkgName}}: no {{$varNameSingular}}Query provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(ctx)
	if err != nil {
	return errors.Wrap(err, "{{.PkgName}}: unable to delete all from {{.Model.Name}}")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o {{$modelNameSingular}}Slice) DeleteAll(ctx context.Context) error {
	if o == nil {
		return errors.New("{{.PkgName}}: no {{$modelNameSingular}} slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	{{ hook . "before_delete_slice" "o" .Model }}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$varNameSingular}}PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM {{$schemaModel}} WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), {{if .Dialect.IndexPlaceholders}}1{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns, len(o))

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to delete all from {{$varNameSingular}} slice")
	}

	{{ hook . "after_delete_slice" "o" .Model }}

	return nil
}
