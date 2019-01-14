{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
// One returns a single {{$varNameSingular}} record from the query. If the query returns no objects, ErrNoRows is returned. 
// If the query returns multiple rows, bunny.ErrMultipleRows is returned.
func (q {{$varNameSingular}}Query) One(ctx context.Context) (*{{$modelNameSingular}}, error) {
	o := &{{$modelNameSingular}}{}

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Wrap(err, "{{.PkgName}}: failed to execute a one query for {{.Model.Name}}")
	}

	{{ hook . "after_select" "o" .Model }}

	return o, nil
}

// First returns a single {{$varNameSingular}} record from the query. If the query returns no objects, ErrNoRows is returned. 
// If the query returns multiple objects, the first one is picked (and no error is generated).
func (q {{$varNameSingular}}Query) First(ctx context.Context) (*{{$modelNameSingular}}, error) {
	o := &{{$modelNameSingular}}{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, o)
	if err != nil {
		return nil, errors.Wrap(err, "{{.PkgName}}: failed to execute a one query for {{.Model.Name}}")
	}

	{{ hook . "after_select" "o" .Model }}

	return o, nil
}

// All returns all {{$modelNameSingular}} records from the query.
func (q {{$varNameSingular}}Query) All(ctx context.Context) ({{$modelNameSingular}}Slice, error) {
	var o []*{{$modelNameSingular}}

	err := q.Bind(ctx, &o)
	if err != nil {
		return nil, errors.Wrap(err, "{{.PkgName}}: failed to assign all query results to {{$modelNameSingular}} slice")
	}

	{{ hook . "after_select_slice" "o" .Model }}

	return o, nil
}

// Count returns the count of all {{$modelNameSingular}} records in the query.
func (q {{$varNameSingular}}Query) Count(ctx context.Context) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to count {{.Model.Name}} rows")
	}

	return count, nil
}

// Exists checks if the row exists in the model.
func (q {{$varNameSingular}}Query) Exists(ctx context.Context) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(ctx).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "{{.PkgName}}: failed to check if {{.Model.Name}} exists")
	}

	return count > 0, nil
}
