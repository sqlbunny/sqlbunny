{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $schemaModel := .Model.Name | schemaModel}}
// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those fields supplied are inserted
// No whitelist behavior: Without a whitelist, fields are inferred by the following rules:
// - All fields without a default value are included (i.e. name, age)
// - All fields with a default, but non-zero are included (i.e. health = 75)
func (o *{{$modelNameSingular}}) Insert(ctx context.Context, whitelist ... {{$modelNameSingular}}Column) error {
	if o == nil {
		return errors.New("{{.PkgName}}: no {{.Model.Name}} provided for insertion")
	}
    _, err := o.InsertIgnore(ctx, "", whitelist...)
	return err
}


func (o *{{$modelNameSingular}}) InsertIgnore(ctx context.Context, ignoreConflictCondition string, whitelist ... {{$modelNameSingular}}Column) (bool, error) {
	if o == nil {
		return false, errors.New("{{.PkgName}}: no {{.Model.Name}} provided for insertion")
	}

	var err error

	{{ hook . "before_insert" "o" .Model }}

	var wl []string
	if len(whitelist) == 0 {
		wl = {{$varNameSingular}}Columns
	} else {
		wl = columnStrings(whitelist)
	}

	key := makeCacheKey(append(wl, ignoreConflictCondition))
	{{$varNameSingular}}InsertCacheMut.RLock()
	cache, cached := {{$varNameSingular}}InsertCache[key]
	{{$varNameSingular}}InsertCacheMut.RUnlock()

	if !cached {
		cache.valueMapping, err = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, wl)
		if err != nil {
			return false, err
		}

		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO {{$schemaModel}} ({{.LQ}}%s{{.RQ}}) VALUES (%s)", strings.Join(wl, "{{.RQ}},{{.LQ}}"), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO {{$schemaModel}} DEFAULT VALUES"
		}

        if len(ignoreConflictCondition) > 0 {
           cache.query += fmt.Sprintf(" ON CONFLICT %s DO NOTHING", ignoreConflictCondition)
        }
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	res, err := bunny.Exec(ctx, cache.query, vals...)
	if err != nil {
		return false, errors.Errorf("{{.PkgName}}: unable to insert into {{.Model.Name}}: %w", err)
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return false, errors.Errorf("{{.PkgName}}: unable to get rows affected for insert into {{.Model.Name}}: %w", err)
	}
	inserted := aff != 0

	if !cached {
		{{$varNameSingular}}InsertCacheMut.Lock()
		{{$varNameSingular}}InsertCache[key] = cache
		{{$varNameSingular}}InsertCacheMut.Unlock()
	}

	{{ hook . "after_insert" "o" .Model }}

	return inserted, nil
}
