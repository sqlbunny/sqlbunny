{{- if .Model.IsJoinModel -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Model.ToManyRelationships -}}
		{{- $varNameSingular := $dot.Model.Name | singular | camelCase -}}
		{{- $txt := txtsFromToMany $dot.Models $dot.Model . -}}
		{{- $arg := printf "maybe%s" $txt.LocalModel.NameGo -}}
		{{- $schemaForeignModel := .ForeignModel | schemaModel}}
// Load{{$txt.Function.NameGo}} allows an eager lookup of values, cached into the
// loaded structs of the objects.
func ({{$varNameSingular}}L) Load{{$txt.Function.NameGo}}(ctx context.Context, singular bool, {{$arg}} interface{}) error {
	var slice []*{{$txt.LocalModel.NameGo}}
	var object *{{$txt.LocalModel.NameGo}}

	count := 1
	if singular {
		object = {{$arg}}.(*{{$txt.LocalModel.NameGo}})
	} else {
		slice = *{{$arg}}.(*[]*{{$txt.LocalModel.NameGo}})
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &{{$varNameSingular}}R{}
		}
		args[0] = object.{{.Column | titleCase}}
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &{{$varNameSingular}}R{}
			}
			args[i] = obj.{{.Column | titleCase}}
		}
	}

		{{if .ToJoinModel -}}
			{{- $schemaJoinModel := .JoinModel | schemaModel -}}
	query := fmt.Sprintf(
		"select {{id 0 | quotes}}.*, {{id 1 | quotes}}.{{.JoinLocalColumn | quotes}} from {{$schemaForeignModel}} as {{id 0 | quotes}} inner join {{$schemaJoinModel}} as {{id 1 | quotes}} on {{id 0 | quotes}}.{{.ForeignColumn | quotes}} = {{id 1 | quotes}}.{{.JoinForeignColumn | quotes}} where {{id 1 | quotes}}.{{.JoinLocalColumn | quotes}} in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
		{{else -}}
	query := fmt.Sprintf(
		"select * from {{$schemaForeignModel}} where {{.ForeignColumn | quotes}} in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
		{{end -}}

	results, err := bunny.Query(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load {{.ForeignModel}}")
	}
	defer results.Close()

	var resultSlice []*{{$txt.ForeignModel.NameGo}}
	{{if .ToJoinModel -}}
	{{- $foreignModel := getModel $dot.Models .ForeignModel -}}
	{{- $joinModel := getModel $dot.Models .JoinModel -}}
	{{- $localCol := $joinModel.GetColumn .JoinLocalColumn}}
	var localJoinCols []{{goType $localCol.GoType}}
	for results.Next() {
		one := new({{$txt.ForeignModel.NameGo}})
		var localJoinCol {{goType $localCol.GoType}}

		err = results.Scan({{$foreignModel.Columns | columnNames | stringMap $dot.StringFuncs.titleCaseIdentifier | prefixStringSlice "&one." | join ", "}}, &localJoinCol)
		if err = results.Err(); err != nil {
			return errors.Wrap(err, "failed to plebian-bind eager loaded slice {{.ForeignModel}}")
		}

		resultSlice = append(resultSlice, one)
		localJoinCols = append(localJoinCols, localJoinCol)
	}

	if err = results.Err(); err != nil {
		return errors.Wrap(err, "failed to plebian-bind eager loaded slice {{.ForeignModel}}")
	}
	{{else -}}
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice {{.ForeignModel}}")
	}
	{{end}}

	{{ $foreignModel := getModel $dot.Models .ForeignModel }}
	{{ hook $dot "after_select_slice_noreturn" "resultSlice" $foreignModel }}

	if singular {
		object.R.{{$txt.Function.NameGo}} = resultSlice
		return nil
	}

	{{if .ToJoinModel -}}
	for i, foreign := range resultSlice {
		localJoinCol := localJoinCols[i]
		for _, local := range slice {
			{{if $txt.Function.UsesBytes -}}
			if 0 == bytes.Compare(local.{{$txt.Function.LocalAssignment}}, localJoinCol) {
			{{else -}}
			if local.{{$txt.Function.LocalAssignment}} == localJoinCol {
			{{end -}}
				local.R.{{$txt.Function.NameGo}} = append(local.R.{{$txt.Function.NameGo}}, foreign)
				break
			}
		}
	}
	{{else -}}
	for _, foreign := range resultSlice {
		for _, local := range slice {
			{{if $txt.Function.UsesBytes -}}
			if 0 == bytes.Compare(local.{{$txt.Function.LocalAssignment}}, foreign.{{$txt.Function.ForeignAssignment}}) {
			{{else -}}
			if local.{{$txt.Function.LocalAssignment}} == foreign.{{$txt.Function.ForeignAssignment}} {
			{{end -}}
				local.R.{{$txt.Function.NameGo}} = append(local.R.{{$txt.Function.NameGo}}, foreign)
				break
			}
		}
	}
	{{end}}

	return nil
}

{{end -}}{{/* range tomany */}}
{{- end -}}{{/* if IsJoinModel */}}
