{{- if .Model.IsJoinModel -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Model.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $dot.Models $dot.Model . -}}
		{{- $varNameSingular := .Model | singular | camelCase -}}
		{{- $foreignVarNameSingular := .ForeignModel | singular | camelCase -}}
		{{- $foreignPrimaryKeyCols := (getModel $dot.Models .ForeignModel).PrimaryKey.Columns -}}
		{{- $foreignSchemaModel := .ForeignModel | $dot.SchemaModel}}
// Set{{$txt.Function.NameGo}} of the {{.Model | singular}} to the related item.
// Sets o.R.{{$txt.Function.NameGo}} to related.
// Adds o to related.R.{{$txt.Function.ForeignNameGo}}.
func (o *{{$txt.LocalModel.NameGo}}) Set{{$txt.Function.NameGo}}(ctx context.Context, insert bool, related *{{$txt.ForeignModel.NameGo}}) error {
	var err error

	if insert {
		related.{{$txt.Function.ForeignAssignment}} = o.{{$txt.Function.LocalAssignment}}
		{{if .ForeignColumnNullable -}}
		related.{{$txt.ForeignModel.ColumnNameGo}}.Valid = true
		{{- end}}

		if err = related.Insert(ctx); err != nil {
			return errors.Wrap(err, "failed to insert into foreign model")
		}
	} else {
		updateQuery := fmt.Sprintf(
			"UPDATE {{$foreignSchemaModel}} SET %s WHERE %s",
			strmangle.SetParamNames("{{$dot.LQ}}", "{{$dot.RQ}}", {{if $dot.Dialect.IndexPlaceholders}}1{{else}}0{{end}}, []string{{"{"}}"{{.ForeignColumn}}"{{"}"}}),
			strmangle.WhereClause("{{$dot.LQ}}", "{{$dot.RQ}}", {{if $dot.Dialect.IndexPlaceholders}}2{{else}}0{{end}}, {{$foreignVarNameSingular}}PrimaryKeyColumns),
		)
		values := []interface{}{o.{{$txt.LocalModel.ColumnNameGo}}, related.{{$foreignPrimaryKeyCols | stringMap $dot.StringFuncs.titleCaseIdentifier | join ", related."}}{{"}"}}

		if _, err = bunny.Exec(ctx, updateQuery, values...); err != nil {
			return errors.Wrap(err, "failed to update foreign model")
		}

		related.{{$txt.Function.ForeignAssignment}} = o.{{$txt.Function.LocalAssignment}}
		{{if .ForeignColumnNullable -}}
		related.{{$txt.ForeignModel.ColumnNameGo}}.Valid = true
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
func (o *{{$txt.LocalModel.NameGo}}) Remove{{$txt.Function.NameGo}}(ctx context.Context, related *{{$txt.ForeignModel.NameGo}}) error {
	var err error

	related.{{$txt.ForeignModel.ColumnNameGo}}.Valid = false
	if err = related.Update(ctx, "{{.ForeignColumn}}"); err != nil {
		related.{{$txt.ForeignModel.ColumnNameGo}}.Valid = true
		return errors.Wrap(err, "failed to update local model")
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
{{- end -}}{{/* join model */}}
