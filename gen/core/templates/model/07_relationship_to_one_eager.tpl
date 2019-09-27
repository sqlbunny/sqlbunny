{{- if .Model.IsJoinModel -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Model.SingleColumnForeignKeys -}}
		{{- $txt := txtsFromFKey $dot.Models $dot.Model . -}}
		{{- $varNameSingular := $dot.Model.Name | singular | camelCase -}}
		{{- $arg := printf "maybe%s" $txt.LocalModel.NameGo -}}
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
		args[0] = object.{{$txt.LocalModel.ColumnNameGo}}
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &{{$varNameSingular}}R{}
			}
			args[i] = obj.{{$txt.LocalModel.ColumnNameGo}}
		}
	}

	query := fmt.Sprintf(
		"select * from {{.ForeignModel | schemaModel}} where {{.ForeignColumn | quotes}} in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	results, err := bunny.Query(ctx, query, args...)
	if err != nil {
		return errors.Errorf("failed to eager load {{$txt.ForeignModel.NameGo}}: %w", err)
	}
	defer results.Close()

	var resultSlice []*{{$txt.ForeignModel.NameGo}}
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Errorf("failed to bind eager loaded slice {{$txt.ForeignModel.NameGo}}: %w", err)
	}

	{{ $foreignModel := getModel $dot.Models .ForeignModel }}
	{{ hook $dot "after_select_slice_noreturn" "resultSlice" $foreignModel }}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		object.R.{{$txt.Function.NameGo}} = resultSlice[0]
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			{{if $txt.Function.UsesBytes -}}
			if 0 == bytes.Compare(local.{{$txt.Function.LocalAssignment}}, foreign.{{$txt.Function.ForeignAssignment}}) {
			{{else -}}
			if local.{{$txt.Function.LocalAssignment}} == foreign.{{$txt.Function.ForeignAssignment}} {
			{{end -}}
				local.R.{{$txt.Function.NameGo}} = foreign
				break
			}
		}
	}

	return nil
}
{{end -}}{{/* range */}}
{{end}}{{/* join model */}}
