{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $schemaModel := .Model.Name | schemaModel}}
// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *{{$modelNameSingular}}) Upsert(ctx context.Context, updateOnConflict bool, conflictFields []string, updateFields []string, whitelist ...string) error {
	if o == nil {
		return errors.New("{{.PkgName}}: no {{.Model.Name}} provided for upsert")
	}

	{{if not .NoHooks -}}
	if err := o.doBeforeUpsertHooks(ctx); err != nil {
		return err
	}
	{{- end}}

	nzDefaults := queries.NonZeroDefaultSet({{$varNameSingular}}ColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictFields {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range updateFields {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	{{$varNameSingular}}UpsertCacheMut.RLock()
	cache, cached := {{$varNameSingular}}UpsertCache[key]
	{{$varNameSingular}}UpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := strmangle.InsertFieldSet(
			{{$varNameSingular}}Columns,
			{{$varNameSingular}}ColumnsWithDefault,
			{{$varNameSingular}}ColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateFieldSet(
			{{$varNameSingular}}Columns,
			{{$varNameSingular}}PrimaryKeyColumns,
			updateFields,
		)

		if len(update) == 0 {
			return errors.New("{{.PkgName}}: unable to upsert {{.Model.Name}}, could not build update field list")
		}

		conflict := conflictFields
		if len(conflict) == 0 {
			conflict = make([]string, len({{$varNameSingular}}PrimaryKeyColumns))
			copy(conflict, {{$varNameSingular}}PrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "{{$schemaModel}}", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	{{if .UseLastInsertID -}}
	{{- $canLastInsertID := .Model.CanLastInsertID -}}
	{{if $canLastInsertID -}}
	result, err := bunny.Exec(ctx, cache.query, vals...)
	{{else -}}
	_, err = bunny.Exec(ctx, cache.query, vals...)
	{{- end}}
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to upsert for {{.Model.Name}}")
	}

	{{if $canLastInsertID -}}
	var lastID int64
	{{- end}}
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	{{if $canLastInsertID -}}
	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	{{$colName := index .Model.PrimaryKey.Columns 0 -}}
	{{- $col := .Model.GetColumn $colName -}}
	{{- $colTitled := $colName | titleCase}}
	o.{{$colTitled}} = {{$col.Type}}(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == {{$varNameSingular}}Mapping["{{$colTitled}}"] {
		goto CacheNoHooks
	}
	{{- end}}

	identifierCols = []interface{}{
		{{range .Model.PrimaryKey.Columns -}}
		o.{{. | titleCaseIdentifier}},
		{{end -}}
	}

	err = bunny.QueryRow(ctx, cache.retQuery, identifierCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to populate default values for {{.Model.Name}}")
	}
	{{- else}}
	if len(cache.retMapping) != 0 {
		err = bunny.QueryRow(ctx, cache.query, vals...).Scan(returns...)
		if bunny.IsErrNoRows(err) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = bunny.Exec(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to upsert {{.Model.Name}}")
	}
	{{- end}}

{{if .UseLastInsertID -}}
CacheNoHooks:
{{end -}}
	if !cached {
		{{$varNameSingular}}UpsertCacheMut.Lock()
		{{$varNameSingular}}UpsertCache[key] = cache
		{{$varNameSingular}}UpsertCacheMut.Unlock()
	}

	{{if not .NoHooks -}}
	return o.doAfterUpsertHooks(ctx)
	{{- else -}}
	return nil
	{{- end}}
}
