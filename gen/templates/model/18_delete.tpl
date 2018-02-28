{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $schemaModel := .Model.Name | .SchemaModel}}
// Delete deletes a single {{$modelNameSingular}} record with an executor.
// Delete will match against the primary key field to find the record to delete.
func (o *{{$modelNameSingular}}) Delete(ctx context.Context) error {
	if o == nil {
	return errors.New("{{.PkgName}}: no {{$modelNameSingular}} provided for delete")
	}

	{{if not .NoHooks -}}
	if err := o.doBeforeDeleteHooks(ctx); err != nil {
	return err
	}
	{{- end}}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), {{$varNameSingular}}PrimaryKeyMapping)
	sql := "DELETE FROM {{$schemaModel}} WHERE {{if .Dialect.IndexPlaceholders}}{{whereClause .LQ .RQ 1 .Model.PrimaryKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Model.PrimaryKey.Columns}}{{end}}"

	_, err := boil.Exec(ctx, sql, args...)
	if err != nil {
	return errors.Wrap(err, "{{.PkgName}}: unable to delete from {{.Model.Name}}")
	}

	{{if not .NoHooks -}}
	if err := o.doAfterDeleteHooks(ctx); err != nil {
	return err
	}
	{{- end}}

	return nil
}

// DeleteAll deletes all matching rows.
func (q {{$varNameSingular}}Query) DeleteAll() error {
	if q.Query == nil {
	return errors.New("{{.PkgName}}: no {{$varNameSingular}}Query provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
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

	{{if not .NoHooks -}}
	if len({{$varNameSingular}}BeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx); err != nil {
				return err
			}
		}
	}
	{{- end}}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$varNameSingular}}PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM {{$schemaModel}} WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), {{if .Dialect.IndexPlaceholders}}1{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns, len(o))

	_, err := boil.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to delete all from {{$varNameSingular}} slice")
	}

	{{if not .NoHooks -}}
	if len({{$varNameSingular}}AfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx); err != nil {
				return err
			}
		}
	}
	{{- end}}

	return nil
}
