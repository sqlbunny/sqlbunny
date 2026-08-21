{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $schemaModel := .Model.Name | schemaModel -}}
{{- $model := .Model}}
// Delete deletes a single {{$modelNameSingular}} record with an executor.
// Delete will match against the primary key field to find the record to delete.
func (o *{{$modelNameSingular}}) Delete(ctx context.Context) error {
	if o == nil {
	return errors.New("{{.PkgName}}: no {{$modelNameSingular}} provided for delete")
	}

	{{ hook . "before_delete" "o" .Model }}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), {{$varNameSingular}}PrimaryKeyMapping)
	sql := "DELETE FROM {{$schemaModel}} WHERE {{whereClause 1 .Model.PrimaryKey.Fields}}"

	_, err := bunny.Exec(ctx, sql, args...)
	if err != nil {
	return errors.Errorf("{{.PkgName}}: unable to delete from {{.Model.Name}}: %w", err)
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
	return errors.Errorf("{{.PkgName}}: unable to delete all from {{.Model.Name}}: %w", err)
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

	// One array parameter per key column, so the parameter count does not
	// grow with the slice.
	{{- range $i, $p := .Model.PrimaryKey.Fields }}
	{{- $f := $model.FindField $p }}
	keys{{$i}} := make([]{{ goType $f.GoType }}, len(o))
	{{- end }}
	for i, obj := range o {
		{{- range $i, $p := .Model.PrimaryKey.Fields }}
		keys{{$i}}[i] = obj.{{ $p | titleCasePath }}
		{{- end }}
	}

	sql := "DELETE FROM {{$schemaModel}} WHERE {{pkArrayWhere 1 .Model .Model.PrimaryKey.Fields}}"

	_, err := bunny.Exec(ctx, sql, {{ range $i, $p := .Model.PrimaryKey.Fields }}{{if $i}}, {{end}}queries.Array(keys{{$i}}){{ end }})
	if err != nil {
		return errors.Errorf("{{.PkgName}}: unable to delete all from {{$varNameSingular}} slice: %w", err)
	}

	{{ hook . "after_delete_slice" "o" .Model }}

	return nil
}
