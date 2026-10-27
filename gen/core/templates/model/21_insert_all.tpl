{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $schemaModel := .Model.Name | schemaModel -}}
{{- $model := .Model -}}

// InsertAll inserts every row in the slice with one INSERT ... SELECT FROM
// unnest(...) statement.
// Whitelist semantics match the per-row Insert.
func (o {{$modelNameSingular}}Slice) InsertAll(ctx context.Context, whitelist ...{{$modelNameSingular}}Column) error {
	return o.bulkInsertImpl(ctx, "", whitelist)
}

// UpsertAll bulk-inserts every row, appending `ON CONFLICT <conflictTarget>
// DO UPDATE SET ...` to every statement. Semantics for conflictTarget and
// updates match the per-row Upsert.
func (o {{$modelNameSingular}}Slice) UpsertAll(ctx context.Context, conflictTarget string, updates map[{{$modelNameSingular}}Column]string, whitelist ...{{$modelNameSingular}}Column) error {
	clause, err := {{$varNameSingular}}UpsertClause(conflictTarget, updates)
	if err != nil {
		return err
	}
	return o.bulkInsertImpl(ctx, clause, whitelist)
}

// bulkInsertImpl runs the bulk INSERT, optionally appending conflictClause
// (the full " ON CONFLICT ..." suffix, or "" for none). Does not RETURN.
func (o {{$modelNameSingular}}Slice) bulkInsertImpl(ctx context.Context, conflictClause string, whitelist []{{$modelNameSingular}}Column) error {
	if len(o) == 0 {
		return nil
	}
	wl, mapping, err := {{$varNameSingular}}BulkInsertPrep(whitelist)
	if err != nil {
		return err
	}
	colsClause := "\"" + strings.Join(wl, "\",\"") + "\""

	src, args := {{$varNameSingular}}BulkValues(o, wl, mapping)
	query := "INSERT INTO {{$schemaModel}} (" + colsClause + ") " + src + conflictClause
	if _, err := bunny.Exec(ctx, query, args...); err != nil {
		return errors.Errorf("{{.PkgName}}: unable to bulk insert into {{.Model.Name}}: %w", err)
	}
	return nil
}

// {{$varNameSingular}}BulkValues builds the "SELECT ... FROM unnest(...)"
// source for a bulk insert of o, with one array parameter per column.
func {{$varNameSingular}}BulkValues(o {{$modelNameSingular}}Slice, wl []string, mapping []queries.MappedField) (string, []any) {
	colTypes := make([]string, len(wl))
	for i, col := range wl {
		colTypes[i] = {{$varNameSingular}}ColumnSQLTypes[col]
	}
	rows := make([][]any, len(o))
	for i, row := range o {
		rows[i] = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(row)), mapping)
	}
	return queries.BulkValues(colTypes, rows)
}

func {{$varNameSingular}}BulkInsertPrep(whitelist []{{$modelNameSingular}}Column) ([]string, []queries.MappedField, error) {
	var wl []string
	if len(whitelist) == 0 {
		wl = {{$varNameSingular}}Columns
	} else {
		wl = columnStrings(whitelist)
	}
	mapping, err := queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, wl)
	if err != nil {
		return nil, nil, err
	}
	return wl, mapping, nil
}

{{if eq (len .Model.PrimaryKey.Fields) 1 -}}
{{- $pkPath := index .Model.PrimaryKey.Fields 0 -}}
{{- $pkColumn := $pkPath.SQLName -}}
{{- $pkField := $model.FindField $pkPath -}}
{{- $pkGoType := goType $pkField.Type.GoType -}}
{{- $pkKeyType := mapKeyType $pkField.Type.GoType }}

// InsertAllIgnore is like InsertAll but appends
// `ON CONFLICT <conflictTarget> DO NOTHING` to every statement. If
// conflictTarget is "" the clause becomes `ON CONFLICT DO NOTHING`.
// Returns the set of primary-key values that were actually inserted; rows
// skipped by the ON CONFLICT clause are absent. Returns nil for an empty
// slice.
{{- if ne $pkKeyType $pkGoType }}
//
// The primary key type {{$pkGoType}} cannot be a map key, so the set has
// {{ mapKey $pkField.Type.GoType "pk" }} keys instead.
{{- end }}
func (o {{$modelNameSingular}}Slice) InsertAllIgnore(ctx context.Context, conflictTarget string, whitelist ...{{$modelNameSingular}}Column) (map[{{$pkKeyType}}]struct{}, error) {
	if len(o) == 0 {
		return nil, nil
	}
	wl, mapping, err := {{$varNameSingular}}BulkInsertPrep(whitelist)
	if err != nil {
		return nil, err
	}
	colsClause := "\"" + strings.Join(wl, "\",\"") + "\""
	conflictClause := {{$varNameSingular}}IgnoreClause(conflictTarget)

	src, args := {{$varNameSingular}}BulkValues(o, wl, mapping)
	query := "INSERT INTO {{$schemaModel}} (" + colsClause + ") " + src + conflictClause + " RETURNING \"{{$pkColumn}}\""

	rows, err := bunny.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Errorf("{{.PkgName}}: unable to bulk insert into {{.Model.Name}}: %w", err)
	}
	defer rows.Close()
	inserted := make(map[{{$pkKeyType}}]struct{}, len(o))
	for rows.Next() {
		var id {{$pkGoType}}
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		inserted[{{ mapKey $pkField.Type.GoType "id" }}] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return inserted, nil
}
{{- end}}
