{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.FKeys -}}
		{{- $txt := txtsFromFKey $dot.Tables $dot.Table . -}}
		{{- $foreignNameSingular := .ForeignTable | singular | camelCase -}}
		{{- $varNameSingular := .Table | singular | camelCase}}
		{{- $schemaTable := .Table | $dot.SchemaTable}}

// Set{{$txt.Function.NameGo}} of the {{.Table | singular}} to the related item.
// Sets o.R.{{$txt.Function.NameGo}} to related.
// Adds o to related.R.{{$txt.Function.ForeignNameGo}}.
func (o *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.NameGo}}(ctx context.Context, insert bool, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error
	if insert {
		if err = related.Insert(ctx); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE {{$schemaTable}} SET %s WHERE %s",
		strmangle.SetParamNames("{{$dot.LQ}}", "{{$dot.RQ}}", {{if $dot.Dialect.IndexPlaceholders}}1{{else}}0{{end}}, []string{{"{"}}"{{.Column}}"{{"}"}}),
		strmangle.WhereClause("{{$dot.LQ}}", "{{$dot.RQ}}", {{if $dot.Dialect.IndexPlaceholders}}2{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns),
	)
	values := []interface{}{related.{{$txt.ForeignTable.ColumnNameGo}}, o.{{$dot.Table.PKey.Columns | stringMap $dot.StringFuncs.titleCase | join ", o."}}{{"}"}}

	if _, err = boil.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.{{$txt.Function.LocalAssignment}} = related.{{$txt.Function.ForeignAssignment}}
	{{if .Nullable -}}
	o.{{$txt.LocalTable.ColumnNameGo}}.Valid = true
	{{- end}}

	if o.R == nil {
		o.R = &{{$varNameSingular}}R{
			{{$txt.Function.NameGo}}: related,
		}
	} else {
		o.R.{{$txt.Function.NameGo}} = related
	}

	{{if .Unique -}}
	if related.R == nil {
		related.R = &{{$foreignNameSingular}}R{
			{{$txt.Function.ForeignNameGo}}: o,
		}
	} else {
		related.R.{{$txt.Function.ForeignNameGo}} = o
	}
	{{else -}}
	if related.R == nil {
		related.R = &{{$foreignNameSingular}}R{
			{{$txt.Function.ForeignNameGo}}: {{$txt.LocalTable.NameGo}}Slice{{"{"}}o{{"}"}},
		}
	} else {
		related.R.{{$txt.Function.ForeignNameGo}} = append(related.R.{{$txt.Function.ForeignNameGo}}, o)
	}
	{{- end}}

	return nil
}

		{{- if .Nullable}}
// Remove{{$txt.Function.NameGo}} relationship.
// Sets o.R.{{$txt.Function.NameGo}} to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.NameGo}}(ctx context.Context, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error

	o.{{$txt.LocalTable.ColumnNameGo}}.Valid = false
	if err = o.Update(ctx, "{{.Column}}"); err != nil {
		o.{{$txt.LocalTable.ColumnNameGo}}.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.{{$txt.Function.NameGo}} = nil
	if related == nil || related.R == nil {
		return nil
	}

	{{if .Unique -}}
	related.R.{{$txt.Function.ForeignNameGo}} = nil
	{{else -}}
	for i, ri := range related.R.{{$txt.Function.ForeignNameGo}} {
		{{if $txt.Function.UsesBytes -}}
		if 0 != bytes.Compare(o.{{$txt.Function.LocalAssignment}}, ri.{{$txt.Function.LocalAssignment}}) {
		{{else -}}
		if o.{{$txt.Function.LocalAssignment}} != ri.{{$txt.Function.LocalAssignment}} {
		{{end -}}
			continue
		}

		ln := len(related.R.{{$txt.Function.ForeignNameGo}})
		if ln > 1 && i < ln-1 {
			related.R.{{$txt.Function.ForeignNameGo}}[i] = related.R.{{$txt.Function.ForeignNameGo}}[ln-1]
		}
		related.R.{{$txt.Function.ForeignNameGo}} = related.R.{{$txt.Function.ForeignNameGo}}[:ln-1]
		break
	}
	{{end -}}

	return nil
}
{{end -}}{{/* if foreignkey nullable */}}
{{- end -}}{{/* range */}}
{{- end -}}{{/* join table */}}
