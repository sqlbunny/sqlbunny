{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $modelNamePlural := .Model.Name | plural | titleCase -}}
{{- $varNamePlural := .Model.Name | plural | camelCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
func test{{$modelNamePlural}}Find(t *testing.T) {
	t.Parallel()

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

	{{$varNameSingular}}Found, err := Find{{$modelNameSingular}}(tx, {{.Model.PrimaryKey.Columns | stringMap .StringFuncs.titleCaseIdentifier | prefixStringSlice (printf "%s." $varNameSingular) | join ", "}})
	if err != nil {
		t.Error(err)
	}

	if {{$varNameSingular}}Found == nil {
		t.Error("want a record, got nil")
	}
}
