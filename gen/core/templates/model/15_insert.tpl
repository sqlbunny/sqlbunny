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
	_, err := o.insertImpl(ctx, "", whitelist)
	return err
}

// InsertIgnore inserts the row, appending `ON CONFLICT <conflictTarget> DO NOTHING`.
// If conflictTarget is "" the clause becomes `ON CONFLICT DO NOTHING` (matches on
// any unique constraint). Returns whether a row was actually inserted.
func (o *{{$modelNameSingular}}) InsertIgnore(ctx context.Context, conflictTarget string, whitelist ... {{$modelNameSingular}}Column) (bool, error) {
	if o == nil {
		return false, errors.New("{{.PkgName}}: no {{.Model.Name}} provided for insertion")
	}
	return o.insertImpl(ctx, {{$varNameSingular}}IgnoreClause(conflictTarget), whitelist)
}

// Upsert inserts the row, appending `ON CONFLICT <conflictTarget> DO UPDATE SET ...`.
// If conflictTarget is "" the clause becomes `ON CONFLICT DO UPDATE SET ...`.
// updates maps column name to SQL expression for the SET list (e.g. EXCLUDED."foo"
// or "foo" + 1). Must be non-empty.
func (o *{{$modelNameSingular}}) Upsert(ctx context.Context, conflictTarget string, updates map[{{$modelNameSingular}}Column]string, whitelist ... {{$modelNameSingular}}Column) error {
	if o == nil {
		return errors.New("{{.PkgName}}: no {{.Model.Name}} provided for insertion")
	}
	clause, err := {{$varNameSingular}}UpsertClause(conflictTarget, updates)
	if err != nil {
		return err
	}
	_, err = o.insertImpl(ctx, clause, whitelist)
	return err
}

// insertImpl runs INSERT and optionally appends conflictClause (the full
// " ON CONFLICT ..." suffix, or "" for none). Returns whether a row was inserted.
func (o *{{$modelNameSingular}}) insertImpl(ctx context.Context, conflictClause string, whitelist []{{$modelNameSingular}}Column) (bool, error) {
	var err error

	{{ hook . "before_insert" "o" .Model }}

	var wl []string
	if len(whitelist) == 0 {
		wl = {{$varNameSingular}}Columns
	} else {
		wl = columnStrings(whitelist)
	}

	key := makeCacheKey(append(wl, conflictClause))
	{{$varNameSingular}}InsertCacheMut.RLock()
	cache, cached := {{$varNameSingular}}InsertCache[key]
	{{$varNameSingular}}InsertCacheMut.RUnlock()

	if !cached {
		cache.valueMapping, err = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, wl)
		if err != nil {
			return false, err
		}

		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO {{$schemaModel}} ({{.LQ}}%s{{.RQ}}) VALUES (%s)", strings.Join(wl, "{{.RQ}},{{.LQ}}"), paramList(len(wl)))
		} else {
			cache.query = "INSERT INTO {{$schemaModel}} DEFAULT VALUES"
		}

		cache.query += conflictClause
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

// {{$varNameSingular}}IgnoreClause builds the " ON CONFLICT [target] DO NOTHING" suffix.
func {{$varNameSingular}}IgnoreClause(conflictTarget string) string {
	if conflictTarget == "" {
		return " ON CONFLICT DO NOTHING"
	}
	return " ON CONFLICT " + conflictTarget + " DO NOTHING"
}

// {{$varNameSingular}}UpsertClause builds the " ON CONFLICT [target] DO UPDATE SET ..." suffix.
// Returns an error if updates is empty.
func {{$varNameSingular}}UpsertClause(conflictTarget string, updates map[{{$modelNameSingular}}Column]string) (string, error) {
	if len(updates) == 0 {
		return "", errors.New("{{.PkgName}}: Upsert requires at least one update column")
	}
	cols := make([]string, 0, len(updates))
	for c := range updates {
		cols = append(cols, string(c))
	}
	sort.Strings(cols)
	var sb strings.Builder
	sb.WriteString(" ON CONFLICT ")
	if conflictTarget != "" {
		sb.WriteString(conflictTarget)
		sb.WriteString(" ")
	}
	sb.WriteString("DO UPDATE SET ")
	for i, c := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("{{.LQ}}")
		sb.WriteString(c)
		sb.WriteString("{{.RQ}} = ")
		sb.WriteString(updates[{{$modelNameSingular}}Column(c)])
	}
	return sb.String(), nil
}
