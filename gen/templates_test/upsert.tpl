{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $modelNamePlural := .Model.Name | plural | titleCase -}}
{{- $varNamePlural := .Model.Name | plural | camelCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
func test{{$modelNamePlural}}Upsert(t *testing.T) {
	t.Parallel()

	if len({{$varNameSingular}}Columns) == len({{$varNameSingular}}PrimaryKeyColumns) {
		t.Skip("Skipping model with only primary key fields")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	{{$varNameSingular}} := {{$modelNameSingular}}{}
	if err = randomize.Struct(seed, &{{$varNameSingular}}, {{$varNameSingular}}DBTypes, true); err != nil {
		t.Errorf("Unable to randomize {{$modelNameSingular}} struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = {{$varNameSingular}}.Upsert(tx, {{if eq .DriverName "postgres"}}false, nil, {{end}}nil); err != nil {
		t.Errorf("Unable to upsert {{$modelNameSingular}}: %s", err)
	}

	count, err := {{$modelNamePlural}}(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &{{$varNameSingular}}, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}PrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize {{$modelNameSingular}} struct: %s", err)
	}

	if err = {{$varNameSingular}}.Upsert(tx, {{if eq .DriverName "postgres"}}true, nil, {{end}}nil); err != nil {
		t.Errorf("Unable to upsert {{$modelNameSingular}}: %s", err)
	}

	count, err = {{$modelNamePlural}}(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
