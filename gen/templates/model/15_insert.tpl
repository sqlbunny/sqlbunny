{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $schemaModel := .Model.Name | .SchemaModel}}
// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those fields supplied are inserted
// No whitelist behavior: Without a whitelist, fields are inferred by the following rules:
// - All fields without a default value are included (i.e. name, age)
// - All fields with a default, but non-zero are included (i.e. health = 75)
func (o *{{$modelNameSingular}}) Insert(ctx context.Context, whitelist ... string) error {
	if o == nil {
		return errors.New("{{.PkgName}}: no {{.Model.Name}} provided for insertion")
	}

	var err error

	{{if not .NoHooks -}}
	if err := o.doBeforeInsertHooks(ctx); err != nil {
		return err
	}
	{{- end}}

	nzDefaults := queries.NonZeroDefaultSet({{$varNameSingular}}ColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	{{$varNameSingular}}InsertCacheMut.RLock()
	cache, cached := {{$varNameSingular}}InsertCache[key]
	{{$varNameSingular}}InsertCacheMut.RUnlock()

	if !cached {
		wl, returnFields := strmangle.InsertFieldSet(
			{{$varNameSingular}}Columns,
			{{$varNameSingular}}ColumnsWithDefault,
			{{$varNameSingular}}ColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, returnFields)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO {{$schemaModel}} ({{.LQ}}%s{{.RQ}}) %%sVALUES (%s)%%s", strings.Join(wl, "{{.RQ}},{{.LQ}}"), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO {{$schemaModel}} DEFAULT VALUES"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			{{if .UseLastInsertID -}}
			cache.retQuery = fmt.Sprintf("SELECT {{.LQ}}%s{{.RQ}} FROM {{$schemaModel}} WHERE %s", strings.Join(returnFields, "{{.RQ}},{{.LQ}}"), strmangle.WhereClause("{{.LQ}}", "{{.RQ}}", {{if .Dialect.IndexPlaceholders}}1{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns))
			{{else -}}
			queryReturning = fmt.Sprintf(" RETURNING {{.LQ}}%s{{.RQ}}", strings.Join(returnFields, "{{.RQ}},{{.LQ}}"))
			{{end -}}
		}

		if len(wl) != 0 {
			cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	{{if .UseLastInsertID -}}
	{{- $canLastInsertID := .Model.CanLastInsertID -}}
	{{if $canLastInsertID -}}
	result, err := boil.Exec(ctx, cache.query, vals...)
	{{else -}}
	_, err = boil.Exec(ctx, cache.query, vals...)
	{{- end}}
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to insert into {{.Model.Name}}")
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
	{{- $col := .Model.GetField $colName -}}
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

	err = boil.QueryRow(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to populate default values for {{.Model.Name}}")
	}
	{{else}}
	if len(cache.retMapping) != 0 {
		err = boil.QueryRow(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = boil.Exec(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to insert into {{.Model.Name}}")
	}
	{{end}}

{{if .UseLastInsertID -}}
CacheNoHooks:
{{- end}}
	if !cached {
		{{$varNameSingular}}InsertCacheMut.Lock()
		{{$varNameSingular}}InsertCache[key] = cache
		{{$varNameSingular}}InsertCacheMut.Unlock()
	}

	{{if not .NoHooks -}}
	return o.doAfterInsertHooks(ctx)
	{{- else -}}
	return nil
	{{- end}}
}
