{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $modelNamePlural := .Model.Name | plural | titleCase -}}
{{- $varNamePlural := .Model.Name | plural | camelCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
func test{{$modelNamePlural}}Bind(t *testing.T) {
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

	if err = {{$modelNamePlural}}(tx).Bind({{$varNameSingular}}); err != nil {
		t.Error(err)
	}
}

func test{{$modelNamePlural}}One(t *testing.T) {
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

	if x, err := {{$modelNamePlural}}(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func test{{$modelNamePlural}}All(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	{{$varNameSingular}}One := &{{$modelNameSingular}}{}
	{{$varNameSingular}}Two := &{{$modelNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}One, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$modelNameSingular}} struct: %s", err)
	}
	if err = randomize.Struct(seed, {{$varNameSingular}}Two, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$modelNameSingular}} struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = {{$varNameSingular}}One.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = {{$varNameSingular}}Two.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := {{$modelNamePlural}}(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func test{{$modelNamePlural}}Count(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	{{$varNameSingular}}One := &{{$modelNameSingular}}{}
	{{$varNameSingular}}Two := &{{$modelNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}One, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$modelNameSingular}} struct: %s", err)
	}
	if err = randomize.Struct(seed, {{$varNameSingular}}Two, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$modelNameSingular}} struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = {{$varNameSingular}}One.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = {{$varNameSingular}}Two.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := {{$modelNamePlural}}(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
