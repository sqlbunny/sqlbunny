{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $modelNameSingularCamel := .Model.Name | singular | camelCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $varNamePlural := .Model.Name | plural | camelCase -}}
{{- $schemaModel := .Model.Name | schemaModel}}
// Reload refetches the object from the database
// using the primary keys with an executor.
//
// Any relationships that were loaded via qm.Load are re-loaded against the
// refreshed row so they stay consistent with the new field values.
func (o *{{$modelNameSingular}}) Reload(ctx context.Context) error {
	{{range .Model.Relationships -}}
	{{ .Name | camelCase }}Loaded := o.R.{{ .Name | camelCase }}Loaded
	{{end}}

	ret, err := Find{{$modelNameSingular}}(ctx {{range .Model.PrimaryKey.Fields}}, o.{{. | titleCasePath}}{{end}})
	if err != nil {
		return err
	}

	*o = *ret

	{{if .Model.Relationships -}}
	slice := []*{{$modelNameSingular}}{o}
	_ = slice
	{{range .Model.Relationships -}}
	if {{ .Name | camelCase }}Loaded {
		if err := ({{$modelNameSingularCamel}}L{}).Load{{ .Name | titleCase }}(ctx, slice); err != nil {
			return err
		}
	}
	{{end -}}
	{{- end}}

	return nil
}

// ReloadAll refetches every row with matching primary key field values
// and overwrites the original object slice with the newly updated slice.
//
// Any relationships loaded on any element of the slice are re-loaded on the
// full slice after the row refetch.
func (o *{{$modelNameSingular}}Slice) ReloadAll(ctx context.Context) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	{{range .Model.Relationships -}}
	{{ .Name | camelCase }}Loaded := false
	{{end}}
	{{if .Model.Relationships -}}
	for _, obj := range *o {
		{{range .Model.Relationships -}}
		if obj.R.{{ .Name | camelCase }}Loaded {
			{{ .Name | camelCase }}Loaded = true
		}
		{{end -}}
	}
	{{- end}}

	{{$varNamePlural}} := {{$modelNameSingular}}Slice{}
	var args []any
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$varNameSingular}}PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT {{$schemaModel}}.* FROM {{$schemaModel}} WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), {{if .Dialect.IndexPlaceholders}}1{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, &{{$varNamePlural}})
	if err != nil {
		return errors.Errorf("{{.PkgName}}: unable to reload all in {{$modelNameSingular}}Slice: %w", err)
	}

	*o = {{$varNamePlural}}

	{{range .Model.Relationships -}}
	if {{ .Name | camelCase }}Loaded {
		if err := ({{$modelNameSingularCamel}}L{}).Load{{ .Name | titleCase }}(ctx, *o); err != nil {
			return err
		}
	}
	{{end}}

	return nil
}
