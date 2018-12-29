{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
// One returns a single {{$varNameSingular}} record from the query.
func (q {{$varNameSingular}}Query) One() (*{{$modelNameSingular}}, error) {
	ctx := queries.GetContext(q.Query)
	_ = ctx  // Suppress "ctx declared and not used" errors. ctx may be used by the hook below.

	o := &{{$modelNameSingular}}{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		return nil, errors.Wrap(err, "{{.PkgName}}: failed to execute a one query for {{.Model.Name}}")
	}

	{{ hook . "after_select" "o" .Model }}

	return o, nil
}

// All returns all {{$modelNameSingular}} records from the query.
func (q {{$varNameSingular}}Query) All() ({{$modelNameSingular}}Slice, error) {
	ctx := queries.GetContext(q.Query)
	_ = ctx  // Suppress "ctx declared and not used" errors. ctx may be used by the hook below.

	var o []*{{$modelNameSingular}}

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "{{.PkgName}}: failed to assign all query results to {{$modelNameSingular}} slice")
	}

	{{ hook . "after_select_slice" "o" .Model }}

	return o, nil
}

// Count returns the count of all {{$modelNameSingular}} records in the query.
func (q {{$varNameSingular}}Query) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to count {{.Model.Name}} rows")
	}

	return count, nil
}

// Exists checks if the row exists in the model.
func (q {{$varNameSingular}}Query) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "{{.PkgName}}: failed to check if {{.Model.Name}} exists")
	}

	return count > 0, nil
}
