{{- if .Model.IsJoinModel -}}
{{- else -}}
	{{- $dot := . -}}
	{{- $model := .Model -}}
	{{- range .Model.ToManyRelationships -}}
		{{- $txt := txtsFromToMany $dot.Models $model . -}}
		{{- $varNameSingular := .Model | singular | camelCase -}}
		{{- $foreignVarNameSingular := .ForeignModel | singular | camelCase}}
		{{- $foreignPrimaryKeyCols := (getModel $dot.Models .ForeignModel).PrimaryKey.Columns -}}
		{{- $foreignSchemaModel := .ForeignModel | $dot.SchemaModel}}
// Add{{$txt.Function.NameGo}} adds the given related objects to the existing relationships
// of the {{$model.Name | singular}}, optionally inserting them as new records.
// Appends related to o.R.{{$txt.Function.NameGo}}.
// Sets related.R.{{$txt.Function.ForeignNameGo}} appropriately.
func (o *{{$txt.LocalModel.NameGo}}) Add{{$txt.Function.NameGo}}(ctx context.Context, insert bool, related ...*{{$txt.ForeignModel.NameGo}}) error {
	var err error
	for _, rel := range related {
		if insert {
			{{if not .ToJoinModel -}}
			rel.{{$txt.Function.ForeignAssignment}} = o.{{$txt.Function.LocalAssignment}}
				{{if .ForeignColumnNullable -}}
			rel.{{$txt.ForeignModel.ColumnNameGo}}.Valid = true
				{{end -}}
			{{end -}}

			if err = rel.Insert(ctx); err != nil {
				return errors.Wrap(err, "failed to insert into foreign model")
			}
		}{{if not .ToJoinModel}} else {
			updateQuery := fmt.Sprintf(
				"UPDATE {{$foreignSchemaModel}} SET %s WHERE %s",
				strmangle.SetParamNames("{{$dot.LQ}}", "{{$dot.RQ}}", {{if $dot.Dialect.IndexPlaceholders}}1{{else}}0{{end}}, []string{{"{"}}"{{.ForeignColumn}}"{{"}"}}),
				strmangle.WhereClause("{{$dot.LQ}}", "{{$dot.RQ}}", {{if $dot.Dialect.IndexPlaceholders}}2{{else}}0{{end}}, {{$foreignVarNameSingular}}PrimaryKeyColumns),
			)
			values := []interface{}{o.{{$txt.LocalModel.ColumnNameGo}}, rel.{{$foreignPrimaryKeyCols | stringMap $dot.StringFuncs.titleCaseIdentifier | join ", rel."}}{{"}"}}

			if _, err = boil.Exec(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign model")
			}

			rel.{{$txt.Function.ForeignAssignment}} = o.{{$txt.Function.LocalAssignment}}
			{{if .ForeignColumnNullable -}}
			rel.{{$txt.ForeignModel.ColumnNameGo}}.Valid = true
			{{end -}}
		}{{end -}}
	}

	{{if .ToJoinModel -}}
	for _, rel := range related {
		query := "insert into {{.JoinModel | $dot.SchemaModel}} ({{.JoinLocalField | $dot.Quotes}}, {{.JoinForeignColumn | $dot.Quotes}}) values {{if $dot.Dialect.IndexPlaceholders}}($1, $2){{else}}(?, ?){{end}}"
		values := []interface{}{{"{"}}o.{{$txt.LocalModel.ColumnNameGo}}, rel.{{$txt.ForeignModel.ColumnNameGo}}}

		_, err = boil.Exec(ctx, query, values...)
		if err != nil {
			return errors.Wrap(err, "failed to insert into join model")
		}
	}
	{{end -}}

	if o.R == nil {
		o.R = &{{$varNameSingular}}R{
			{{$txt.Function.NameGo}}: related,
		}
	} else {
		o.R.{{$txt.Function.NameGo}} = append(o.R.{{$txt.Function.NameGo}}, related...)
	}

	{{if .ToJoinModel -}}
	for _, rel := range related {
		if rel.R == nil {
			rel.R = &{{$foreignVarNameSingular}}R{
				{{$txt.Function.ForeignNameGo}}: {{$txt.LocalModel.NameGo}}Slice{{"{"}}o{{"}"}},
			}
		} else {
			rel.R.{{$txt.Function.ForeignNameGo}} = append(rel.R.{{$txt.Function.ForeignNameGo}}, o)
		}
	}
	{{else -}}
	for _, rel := range related {
		if rel.R == nil {
			rel.R = &{{$foreignVarNameSingular}}R{
				{{$txt.Function.ForeignNameGo}}: o,
			}
		} else {
			rel.R.{{$txt.Function.ForeignNameGo}} = o
		}
	}
	{{end -}}

	return nil
}

			{{- if (or .ForeignColumnNullable .ToJoinModel)}}
// Set{{$txt.Function.NameGo}} removes all previously related items of the
// {{$model.Name | singular}} replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.{{$txt.Function.ForeignNameGo}}'s {{$txt.Function.NameGo}} accordingly.
// Replaces o.R.{{$txt.Function.NameGo}} with related.
// Sets related.R.{{$txt.Function.ForeignNameGo}}'s {{$txt.Function.NameGo}} accordingly.
func (o *{{$txt.LocalModel.NameGo}}) Set{{$txt.Function.NameGo}}(ctx context.Context, insert bool, related ...*{{$txt.ForeignModel.NameGo}}) error {
	{{if .ToJoinModel -}}
	query := "delete from {{.JoinModel | $dot.SchemaModel}} where {{.JoinLocalField | $dot.Quotes}} = {{if $dot.Dialect.IndexPlaceholders}}$1{{else}}?{{end}}"
	values := []interface{}{{"{"}}o.{{$txt.LocalModel.ColumnNameGo}}}
	{{else -}}
	query := "update {{.ForeignModel | $dot.SchemaModel}} set {{.ForeignColumn | $dot.Quotes}} = null where {{.ForeignColumn | $dot.Quotes}} = {{if $dot.Dialect.IndexPlaceholders}}$1{{else}}?{{end}}"
	values := []interface{}{{"{"}}o.{{$txt.LocalModel.ColumnNameGo}}}
	{{end -}}
	_, err := boil.Exec(ctx, query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	{{if .ToJoinModel -}}
	remove{{$txt.Function.NameGo}}From{{$txt.Function.ForeignNameGo}}Slice(o, related)
	if o.R != nil {
		o.R.{{$txt.Function.NameGo}} = nil
	}
	{{else -}}
	if o.R != nil {
		for _, rel := range o.R.{{$txt.Function.NameGo}} {
			rel.{{$txt.ForeignModel.ColumnNameGo}}.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.{{$txt.Function.ForeignNameGo}} = nil
		}

		o.R.{{$txt.Function.NameGo}} = nil
	}
	{{end -}}

	return o.Add{{$txt.Function.NameGo}}(ctx, insert, related...)
}

// Remove{{$txt.Function.NameGo}} relationships from objects passed in.
// Removes related items from R.{{$txt.Function.NameGo}} (uses pointer comparison, removal does not keep order)
// Sets related.R.{{$txt.Function.ForeignNameGo}}.
func (o *{{$txt.LocalModel.NameGo}}) Remove{{$txt.Function.NameGo}}(ctx context.Context, related ...*{{$txt.ForeignModel.NameGo}}) error {
	var err error
	{{if .ToJoinModel -}}
	query := fmt.Sprintf(
		"delete from {{.JoinModel | $dot.SchemaModel}} where {{.JoinLocalField | $dot.Quotes}} = {{if $dot.Dialect.IndexPlaceholders}}$1{{else}}?{{end}} and {{.JoinForeignColumn | $dot.Quotes}} in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, len(related), 2, 1),
	)
	values := []interface{}{{"{"}}o.{{$txt.LocalModel.ColumnNameGo}}}
	for _, rel := range related {
		values = append(values, rel.{{$txt.ForeignModel.ColumnNameGo}})
	}

	_, err = boil.Exec(ctx, query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}
	{{else -}}
	for _, rel := range related {
		rel.{{$txt.ForeignModel.ColumnNameGo}}.Valid = false
		{{if not .ToJoinModel -}}
		if rel.R != nil {
			rel.R.{{$txt.Function.ForeignNameGo}} = nil
		}
		{{end -}}
		if err = rel.Update(ctx, "{{.ForeignColumn}}"); err != nil {
			return err
		}
	}
	{{end -}}

	{{if .ToJoinModel -}}
	remove{{$txt.Function.NameGo}}From{{$txt.Function.ForeignNameGo}}Slice(o, related)
	{{end -}}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.{{$txt.Function.NameGo}} {
			if rel != ri {
				continue
			}

			ln := len(o.R.{{$txt.Function.NameGo}})
			if ln > 1 && i < ln-1 {
				o.R.{{$txt.Function.NameGo}}[i] = o.R.{{$txt.Function.NameGo}}[ln-1]
			}
			o.R.{{$txt.Function.NameGo}} = o.R.{{$txt.Function.NameGo}}[:ln-1]
			break
		}
	}

	return nil
}

				{{if .ToJoinModel -}}
func remove{{$txt.Function.NameGo}}From{{$txt.Function.ForeignNameGo}}Slice(o *{{$txt.LocalModel.NameGo}}, related []*{{$txt.ForeignModel.NameGo}}) {
	for _, rel := range related {
		if rel.R == nil {
			continue
		}
		for i, ri := range rel.R.{{$txt.Function.ForeignNameGo}} {
			{{if $txt.Function.UsesBytes -}}
			if 0 != bytes.Compare(o.{{$txt.Function.LocalAssignment}}, ri.{{$txt.Function.LocalAssignment}}) {
			{{else -}}
			if o.{{$txt.Function.LocalAssignment}} != ri.{{$txt.Function.LocalAssignment}} {
			{{end -}}
				continue
			}

			ln := len(rel.R.{{$txt.Function.ForeignNameGo}})
			if ln > 1 && i < ln-1 {
				rel.R.{{$txt.Function.ForeignNameGo}}[i] = rel.R.{{$txt.Function.ForeignNameGo}}[ln-1]
			}
			rel.R.{{$txt.Function.ForeignNameGo}} = rel.R.{{$txt.Function.ForeignNameGo}}[:ln-1]
			break
		}
	}
}
				{{end -}}{{- /* if ToJoinModel */ -}}
			{{- end -}}{{- /* if nullable foreign key */ -}}
	{{- end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* if IsJoinModel */ -}}
