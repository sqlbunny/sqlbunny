{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $schemaModel := .Model.Name | schemaModel -}}
{{- $model := .Model -}}

// InsertAll inserts every row in the slice using one or more multi-row
// INSERT statements (chunked to stay under Postgres' bind-parameter limit).
// Whitelist semantics match the per-row Insert.
func (o {{$modelNameSingular}}Slice) InsertAll(ctx context.Context, whitelist ...{{$modelNameSingular}}Column) error {
	if len(o) == 0 {
		return nil
	}
	wl, mapping, err := {{$varNameSingular}}BulkInsertPrep(whitelist)
	if err != nil {
		return err
	}
	colsClause := "{{.LQ}}" + strings.Join(wl, "{{.RQ}},{{.LQ}}") + "{{.RQ}}"

	const maxParams = 60000
	chunkSize := maxParams / len(wl)
	if chunkSize < 1 {
		chunkSize = 1
	}

	for start := 0; start < len(o); start += chunkSize {
		end := start + chunkSize
		if end > len(o) {
			end = len(o)
		}
		chunk := o[start:end]

		args := make([]any, 0, len(chunk)*len(wl))
		for _, row := range chunk {
			args = append(args, queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(row)), mapping)...)
		}

		placeholders := strmangle.Placeholders(dialect.IndexPlaceholders, len(chunk)*len(wl), 1, len(wl))
		query := "INSERT INTO {{$schemaModel}} (" + colsClause + ") VALUES " + placeholders
		if _, err := bunny.Exec(ctx, query, args...); err != nil {
			return errors.Errorf("{{.PkgName}}: unable to bulk insert into {{.Model.Name}}: %w", err)
		}
	}
	return nil
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
{{- $pkGoType := goType ($model.FindField $pkPath).Type.GoType }}

// InsertAllIgnore is like InsertAll but appends
// `ON CONFLICT <ignoreConflictCondition> DO NOTHING` to every statement.
// Returns the set of primary-key values that were actually inserted; rows
// skipped by the ON CONFLICT clause are absent. Returns nil for an empty
// slice.
func (o {{$modelNameSingular}}Slice) InsertAllIgnore(ctx context.Context, ignoreConflictCondition string, whitelist ...{{$modelNameSingular}}Column) (map[{{$pkGoType}}]struct{}, error) {
	if len(o) == 0 {
		return nil, nil
	}
	wl, mapping, err := {{$varNameSingular}}BulkInsertPrep(whitelist)
	if err != nil {
		return nil, err
	}
	colsClause := "{{.LQ}}" + strings.Join(wl, "{{.RQ}},{{.LQ}}") + "{{.RQ}}"

	const maxParams = 60000
	chunkSize := maxParams / len(wl)
	if chunkSize < 1 {
		chunkSize = 1
	}

	inserted := make(map[{{$pkGoType}}]struct{}, len(o))
	for start := 0; start < len(o); start += chunkSize {
		end := start + chunkSize
		if end > len(o) {
			end = len(o)
		}
		chunk := o[start:end]

		args := make([]any, 0, len(chunk)*len(wl))
		for _, row := range chunk {
			args = append(args, queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(row)), mapping)...)
		}

		placeholders := strmangle.Placeholders(dialect.IndexPlaceholders, len(chunk)*len(wl), 1, len(wl))
		query := "INSERT INTO {{$schemaModel}} (" + colsClause + ") VALUES " + placeholders
		if len(ignoreConflictCondition) > 0 {
			query += " ON CONFLICT " + ignoreConflictCondition + " DO NOTHING"
		}
		query += " RETURNING {{.LQ}}{{$pkColumn}}{{.RQ}}"

		rows, err := bunny.Query(ctx, query, args...)
		if err != nil {
			return nil, errors.Errorf("{{.PkgName}}: unable to bulk insert into {{.Model.Name}}: %w", err)
		}
		for rows.Next() {
			var id {{$pkGoType}}
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return nil, err
			}
			inserted[id] = struct{}{}
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, err
		}
		rows.Close()
	}
	return inserted, nil
}
{{- end}}
