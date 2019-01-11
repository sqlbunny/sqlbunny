{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $schemaModel := .Model.Name | schemaModel}}
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

	{{ hook . "before_insert" "o" .Model }}

	if len(whitelist) == 0 {
		whitelist = {{$varNameSingular}}Columns
	}

	key := makeCacheKey(whitelist)
	{{$varNameSingular}}InsertCacheMut.RLock()
	cache, cached := {{$varNameSingular}}InsertCache[key]
	{{$varNameSingular}}InsertCacheMut.RUnlock()

	if !cached {
		cache.valueMapping, err = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, whitelist)
		if err != nil {
			return err
		}

		if len(whitelist) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO {{$schemaModel}} ({{.LQ}}%s{{.RQ}}) VALUES (%s)", strings.Join(whitelist, "{{.RQ}},{{.LQ}}"), strmangle.Placeholders(dialect.IndexPlaceholders, len(whitelist), 1, 1))
		} else {
			cache.query = "INSERT INTO {{$schemaModel}} DEFAULT VALUES"
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	_, err = bunny.Exec(ctx, cache.query, vals...)
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to insert into {{.Model.Name}}")
	}

	if !cached {
		{{$varNameSingular}}InsertCacheMut.Lock()
		{{$varNameSingular}}InsertCache[key] = cache
		{{$varNameSingular}}InsertCacheMut.Unlock()
	}

	{{ hook . "after_insert" "o" .Model }}
	
	return nil
}
