{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $schemaModel := .Model.Name | .SchemaModel}}
// Update uses an executor to update the {{$modelNameSingular}}.
// Whitelist behavior: If a whitelist is provided, only the fields given are updated.
// No whitelist behavior: Without a whitelist, fields are inferred by the following rules:
// - All fields are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *{{$modelNameSingular}}) Update(ctx context.Context, whitelist ... string) error {
	var err error
	{{if not .NoHooks -}}
	if err = o.doBeforeUpdateHooks(ctx); err != nil {
		return err
	}
	{{end -}}

	key := makeCacheKey(whitelist, nil)
	{{$varNameSingular}}UpdateCacheMut.RLock()
	cache, cached := {{$varNameSingular}}UpdateCache[key]
	{{$varNameSingular}}UpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateFieldSet(
			{{$varNameSingular}}Columns,
			{{$varNameSingular}}PrimaryKeyColumns,
			whitelist,
		)
		if len(wl) == 0 {
			return errors.New("{{.PkgName}}: unable to update {{.Model.Name}}, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE {{$schemaModel}} SET %s WHERE %s",
			strmangle.SetParamNames("{{.LQ}}", "{{.RQ}}", {{if .Dialect.IndexPlaceholders}}1{{else}}0{{end}}, wl),
			strmangle.WhereClause("{{.LQ}}", "{{.RQ}}", {{if .Dialect.IndexPlaceholders}}len(wl)+1{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, append(wl, {{$varNameSingular}}PrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	_, err = boil.Exec(ctx, cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to update {{.Model.Name}} row")
	}

	if !cached {
		{{$varNameSingular}}UpdateCacheMut.Lock()
		{{$varNameSingular}}UpdateCache[key] = cache
		{{$varNameSingular}}UpdateCacheMut.Unlock()
	}

	{{if not .NoHooks -}}
	return o.doAfterUpdateHooks(ctx)
	{{- else -}}
	return nil
	{{- end}}
}

// UpdateAll updates all rows with the specified field values.
func (q {{$varNameSingular}}Query) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to update all for {{.Model.Name}}")
	}

	return nil
}

// UpdateAll updates all rows with the specified field values, using an executor.
func (o {{$modelNameSingular}}Slice) UpdateAll(ctx context.Context, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("{{.PkgName}}: update all requires at least one field argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each field
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$varNameSingular}}PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE {{$schemaModel}} SET %s WHERE %s",
		strmangle.SetParamNames("{{.LQ}}", "{{.RQ}}", {{if .Dialect.IndexPlaceholders}}1{{else}}0{{end}}, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), {{if .Dialect.IndexPlaceholders}}len(colNames)+1{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns, len(o)))

	_, err := boil.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to update all in {{$varNameSingular}} slice")
	}

	return nil
}
