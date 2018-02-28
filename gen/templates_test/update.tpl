{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $modelNamePlural := .Model.Name | plural | titleCase -}}
{{- $varNamePlural := .Model.Name | plural | camelCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
func test{{$modelNamePlural}}Update(t *testing.T) {
	t.Parallel()

	if len({{$varNameSingular}}Columns) == len({{$varNameSingular}}PrimaryKeyColumns) {
		t.Skip("Skipping model with only primary key fields")
	}

	seed := randomize.NewSeed()
	var err error
	{{$varNameSingular}} := &{{$modelNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$modelNameSingular}} struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = {{$varNameSingular}}.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := {{$modelNamePlural}}(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, {{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$modelNameSingular}} struct: %s", err)
	}

	if err = {{$varNameSingular}}.Update(tx); err != nil {
		t.Error(err)
	}
}

func test{{$modelNamePlural}}SliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len({{$varNameSingular}}Columns) == len({{$varNameSingular}}PrimaryKeyColumns) {
		t.Skip("Skipping model with only primary key fields")
	}

	seed := randomize.NewSeed()
	var err error
	{{$varNameSingular}} := &{{$modelNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$modelNameSingular}} struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = {{$varNameSingular}}.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := {{$modelNamePlural}}(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, {{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}PrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize {{$modelNameSingular}} struct: %s", err)
	}

	// Remove Primary keys and unique fields from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch({{$varNameSingular}}Columns, {{$varNameSingular}}PrimaryKeyColumns) {
		fields = {{$varNameSingular}}Columns
	} else {
		fields = strmangle.SetComplement(
			{{$varNameSingular}}Columns,
			{{$varNameSingular}}PrimaryKeyColumns,
		)
		{{- if eq .DriverName "mssql"}}
		fields = strmangle.SetComplement(
			fields,
			{{$varNameSingular}}ColumnsWithAuto,
		)
		{{- end}}
	}

	value := reflect.Indirect(reflect.ValueOf({{$varNameSingular}}))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := {{$modelNameSingular}}Slice{{"{"}}{{$varNameSingular}}{{"}"}}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
