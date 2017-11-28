{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $dot.Tables $dot.Table . -}}
		{{- $varNameSingular := .Table | singular | camelCase -}}
		{{- $foreignVarNameSingular := .ForeignTable | singular | camelCase -}}
		{{- $foreignPKeyCols := (getTable $dot.Tables .ForeignTable).PKey.Columns -}}
		{{- $foreignSchemaTable := .ForeignTable | $dot.SchemaTable}}
// Set{{$txt.Function.NameGo}} of the {{.Table | singular}} to the related item.
// Sets o.R.{{$txt.Function.NameGo}} to related.
// Adds o to related.R.{{$txt.Function.ForeignNameGo}}.
func (o *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.NameGo}}(ctx context.Context, insert bool, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error

	if insert {
		related.{{$txt.Function.ForeignAssignment}} = o.{{$txt.Function.LocalAssignment}}
		{{if .ForeignColumnNullable -}}
		related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = true
		{{- end}}

		if err = related.Insert(ctx); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	} else {
		updateQuery := fmt.Sprintf(
			"UPDATE {{$foreignSchemaTable}} SET %s WHERE %s",
			strmangle.SetParamNames("{{$dot.LQ}}", "{{$dot.RQ}}", {{if $dot.Dialect.IndexPlaceholders}}1{{else}}0{{end}}, []string{{"{"}}"{{.ForeignColumn}}"{{"}"}}),
			strmangle.WhereClause("{{$dot.LQ}}", "{{$dot.RQ}}", {{if $dot.Dialect.IndexPlaceholders}}2{{else}}0{{end}}, {{$foreignVarNameSingular}}PrimaryKeyColumns),
		)
		values := []interface{}{o.{{$txt.LocalTable.ColumnNameGo}}, related.{{$foreignPKeyCols | stringMap $dot.StringFuncs.titleCase | join ", related."}}{{"}"}}

		if _, err = boil.ExecContext(ctx, updateQuery, values...); err != nil {
			return errors.Wrap(err, "failed to update foreign table")
		}

		related.{{$txt.Function.ForeignAssignment}} = o.{{$txt.Function.LocalAssignment}}
		{{if .ForeignColumnNullable -}}
		related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = true
		{{- end}}
	}


	if o.R == nil {
		o.R = &{{$varNameSingular}}R{
			{{$txt.Function.NameGo}}: related,
		}
	} else {
		o.R.{{$txt.Function.NameGo}} = related
	}

	if related.R == nil {
		related.R = &{{$foreignVarNameSingular}}R{
			{{$txt.Function.ForeignNameGo}}: o,
		}
	} else {
		related.R.{{$txt.Function.ForeignNameGo}} = o
	}
	return nil
}

		{{- if .ForeignColumnNullable}}
// Remove{{$txt.Function.NameGo}} relationship.
// Sets o.R.{{$txt.Function.NameGo}} to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.NameGo}}(ctx context.Context, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error

	related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = false
	if err = related.Update(ctx, "{{.ForeignColumn}}"); err != nil {
		related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.{{$txt.Function.NameGo}} = nil
	if related == nil || related.R == nil {
		return nil
	}

	related.R.{{$txt.Function.ForeignNameGo}} = nil
	return nil
}
{{end -}}{{/* if foreignkey nullable */}}
{{- end -}}{{/* range */}}
{{- end -}}{{/* join table */}}
