{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $modelNameSingularCamel := .Model.Name | singular | camelCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $varNamePlural := .Model.Name | plural | camelCase -}}
{{- $schemaModel := .Model.Name | schemaModel -}}
{{- $model := .Model}}
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
	// One array parameter per key column, so the parameter count does not
	// grow with the slice.
	{{- range $i, $p := .Model.PrimaryKey.Fields }}
	{{- $f := $model.FindField $p }}
	keys{{$i}} := make([]{{ goType $f.GoType }}, len(*o))
	{{- end }}
	for i, obj := range *o {
		{{- range $i, $p := .Model.PrimaryKey.Fields }}
		keys{{$i}}[i] = obj.{{ $p | titleCasePath }}
		{{- end }}
	}

	sql := "SELECT {{$schemaModel}}.* FROM {{$schemaModel}} WHERE {{pkArrayWhere 1 .Model .Model.PrimaryKey.Fields}}"

	q := queries.Raw(sql, {{ range $i, $p := .Model.PrimaryKey.Fields }}{{if $i}}, {{end}}queries.Array(keys{{$i}}){{ end }})

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
