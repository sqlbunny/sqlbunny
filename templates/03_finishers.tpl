{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
// One returns a single {{$varNameSingular}} record from the query.
func (q {{$varNameSingular}}Query) One() (*{{$tableNameSingular}}, error) {
	o := &{{$tableNameSingular}}{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "{{.PkgName}}: failed to execute a one query for {{.Table.Name}}")
	}

	{{if not .NoHooks -}}
	if err := o.doAfterSelectHooks(queries.GetContext(q.Query)); err != nil {
		return o, err
	}
	{{- end}}

	return o, nil
}

// All returns all {{$tableNameSingular}} records from the query.
func (q {{$varNameSingular}}Query) All() ({{$tableNameSingular}}Slice, error) {
	var o []*{{$tableNameSingular}}

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "{{.PkgName}}: failed to assign all query results to {{$tableNameSingular}} slice")
	}

	{{if not .NoHooks -}}
	if len({{$varNameSingular}}AfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetContext(q.Query)); err != nil {
				return o, err
			}
		}
	}
	{{- end}}

	return o, nil
}

// Count returns the count of all {{$tableNameSingular}} records in the query.
func (q {{$varNameSingular}}Query) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to count {{.Table.Name}} rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q {{$varNameSingular}}Query) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "{{.PkgName}}: failed to check if {{.Table.Name}} exists")
	}

	return count > 0, nil
}
