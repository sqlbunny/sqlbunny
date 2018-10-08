{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $modelNamePlural := .Model.Name | plural | titleCase -}}
{{- $varNamePlural := .Model.Name | plural | camelCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
func test{{$modelNamePlural}}Exists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	{{$varNameSingular}} := &{{$modelNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$modelNameSingular}} struct: %s", err)
	}

	tx := MustTx(bunny.Begin())
	defer tx.Rollback()
	if err = {{$varNameSingular}}.Insert(tx); err != nil {
		t.Error(err)
	}

	{{$pkeyArgs := .Model.PrimaryKey.Columns | stringMap .StringFuncs.titleCaseIdentifier | prefixStringSlice (printf "%s." $varNameSingular) | join ", " -}}
	e, err := {{$modelNameSingular}}Exists(tx, {{$pkeyArgs}})
	if err != nil {
		t.Errorf("Unable to check if {{$modelNameSingular}} exists: %s", err)
	}
	if !e {
		t.Errorf("Expected {{$modelNameSingular}}ExistsG to return true, but got false.")
	}
}
