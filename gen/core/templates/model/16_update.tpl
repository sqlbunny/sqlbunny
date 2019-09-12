{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $schemaModel := .Model.Name | schemaModel}}
// Update uses an executor to update the {{$modelNameSingular}}.
// Whitelist behavior: If a whitelist is provided, only the fields given are updated.
// No whitelist behavior: Without a whitelist, fields are inferred by the following rules:
// - All fields are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *{{$modelNameSingular}}) Update(ctx context.Context, whitelist ... string) error {
	var err error

	{{ hook . "before_update" "o" .Model }}

	if len(whitelist) == 0 {
		whitelist = {{$varNameSingular}}NonPrimaryKeyColumns
	}

	if len(whitelist) == 0 {
		// Nothing to update
		return nil
	}

	key := makeCacheKey(whitelist)
	{{$varNameSingular}}UpdateCacheMut.RLock()
	cache, cached := {{$varNameSingular}}UpdateCache[key]
	{{$varNameSingular}}UpdateCacheMut.RUnlock()

	if !cached {
		cache.query = fmt.Sprintf("UPDATE {{$schemaModel}} SET %s WHERE %s",
			strmangle.SetParamNames("{{.LQ}}", "{{.RQ}}", {{if .Dialect.IndexPlaceholders}}1{{else}}0{{end}}, whitelist),
			strmangle.WhereClause("{{.LQ}}", "{{.RQ}}", {{if .Dialect.IndexPlaceholders}}len(whitelist)+1{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, append(whitelist, {{$varNameSingular}}PrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to update {{.Model.Name}} row")
	}

	if !cached {
		{{$varNameSingular}}UpdateCacheMut.Lock()
		{{$varNameSingular}}UpdateCache[key] = cache
		{{$varNameSingular}}UpdateCacheMut.Unlock()
	}

	{{ hook . "after_update" "o" .Model }}

	return nil
}

// UpdateMapAll updates all rows with the specified field values.
func (q {{$varNameSingular}}Query) UpdateMapAll(ctx context.Context, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to update all for {{.Model.Name}}")
	}

	return nil
}
