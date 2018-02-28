{{- if .Model.IsJoinModel -}}
{{- else -}}
	{{- $dot := . -}}
	{{- $model := .Model -}}
	{{- range .Model.ToManyRelationships -}}
	{{- $varNameSingular := .Model | singular | camelCase -}}
	{{- $foreignVarNameSingular := .ForeignModel | singular | camelCase -}}
	{{- $txt := txtsFromToMany $dot.Models $model .}}
func test{{$txt.LocalModel.NameGo}}ToManyAddOp{{$txt.Function.NameGo}}(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a {{$txt.LocalModel.NameGo}}
	var b, c, d, e {{$txt.ForeignModel.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, strmangle.SetComplement({{$varNameSingular}}PrimaryKeyColumns, {{$varNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*{{$txt.ForeignModel.NameGo}}{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, {{$foreignVarNameSingular}}DBTypes, false, strmangle.SetComplement({{$foreignVarNameSingular}}PrimaryKeyColumns, {{$foreignVarNameSingular}}FieldsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*{{$txt.ForeignModel.NameGo}}{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.Add{{$txt.Function.NameGo}}(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]
		{{- if .ToJoinModel}}

		if first.R.{{$txt.Function.ForeignNameGo}}[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.{{$txt.Function.ForeignNameGo}}[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		{{- else}}

		{{if $txt.Function.UsesBytes -}}
		if 0 != bytes.Compare(a.{{$txt.Function.LocalAssignment}}, first.{{$txt.Function.ForeignAssignment}}) {
			t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, first.{{$txt.Function.ForeignAssignment}})
		}
		if 0 != bytes.Compare(a.{{$txt.Function.LocalAssignment}}, second.{{$txt.Function.ForeignAssignment}}) {
			t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, second.{{$txt.Function.ForeignAssignment}})
		}
		{{else -}}
		if a.{{$txt.Function.LocalAssignment}} != first.{{$txt.Function.ForeignAssignment}} {
			t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, first.{{$txt.Function.ForeignAssignment}})
		}
		if a.{{$txt.Function.LocalAssignment}} != second.{{$txt.Function.ForeignAssignment}} {
			t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, second.{{$txt.Function.ForeignAssignment}})
		}
		{{- end}}

		if first.R.{{$txt.Function.ForeignNameGo}} != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.{{$txt.Function.ForeignNameGo}} != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		{{- end}}

		if a.R.{{$txt.Function.NameGo}}[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.{{$txt.Function.NameGo}}[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.{{$txt.Function.NameGo}}(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i+1)*2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
{{- if (or .ForeignColumnNullable .ToJoinModel)}}

func test{{$txt.LocalModel.NameGo}}ToManySetOp{{$txt.Function.NameGo}}(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a {{$txt.LocalModel.NameGo}}
	var b, c, d, e {{$txt.ForeignModel.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, strmangle.SetComplement({{$varNameSingular}}PrimaryKeyColumns, {{$varNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*{{$txt.ForeignModel.NameGo}}{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, {{$foreignVarNameSingular}}DBTypes, false, strmangle.SetComplement({{$foreignVarNameSingular}}PrimaryKeyColumns, {{$foreignVarNameSingular}}FieldsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.Set{{$txt.Function.NameGo}}(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.{{$txt.Function.NameGo}}(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.Set{{$txt.Function.NameGo}}(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.{{$txt.Function.NameGo}}(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	{{- if .ToJoinModel}}

	// The following checks cannot be implemented since we have no handle
	// to these when we call Set(). Leaving them here as wishful thinking
	// and to let people know there's dragons.
	//
	// if len(b.R.{{$txt.Function.ForeignNameGo}}) != 0 {
	// 	t.Error("relationship was not removed properly from the slice")
	// }
	// if len(c.R.{{$txt.Function.ForeignNameGo}}) != 0 {
	// 	t.Error("relationship was not removed properly from the slice")
	// }
	if d.R.{{$txt.Function.ForeignNameGo}}[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.{{$txt.Function.ForeignNameGo}}[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	{{- else}}

	if b.{{$txt.ForeignModel.ColumnNameGo}}.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.{{$txt.ForeignModel.ColumnNameGo}}.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	{{if $txt.Function.UsesBytes -}}
	if 0 != bytes.Compare(a.{{$txt.Function.LocalAssignment}}, d.{{$txt.Function.ForeignAssignment}}) {
		t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, d.{{$txt.Function.ForeignAssignment}})
	}
	if 0 != bytes.Compare(a.{{$txt.Function.LocalAssignment}}, e.{{$txt.Function.ForeignAssignment}}) {
		t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, e.{{$txt.Function.ForeignAssignment}})
	}
	{{else -}}
	if a.{{$txt.Function.LocalAssignment}} != d.{{$txt.Function.ForeignAssignment}} {
		t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, d.{{$txt.Function.ForeignAssignment}})
	}
	if a.{{$txt.Function.LocalAssignment}} != e.{{$txt.Function.ForeignAssignment}} {
		t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, e.{{$txt.Function.ForeignAssignment}})
	}
	{{- end}}

	if b.R.{{$txt.Function.ForeignNameGo}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.{{$txt.Function.ForeignNameGo}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.{{$txt.Function.ForeignNameGo}} != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.{{$txt.Function.ForeignNameGo}} != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	{{- end}}

	if a.R.{{$txt.Function.NameGo}}[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.{{$txt.Function.NameGo}}[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func test{{$txt.LocalModel.NameGo}}ToManyRemoveOp{{$txt.Function.NameGo}}(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a {{$txt.LocalModel.NameGo}}
	var b, c, d, e {{$txt.ForeignModel.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, strmangle.SetComplement({{$varNameSingular}}PrimaryKeyColumns, {{$varNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*{{$txt.ForeignModel.NameGo}}{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, {{$foreignVarNameSingular}}DBTypes, false, strmangle.SetComplement({{$foreignVarNameSingular}}PrimaryKeyColumns, {{$foreignVarNameSingular}}FieldsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.Add{{$txt.Function.NameGo}}(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.{{$txt.Function.NameGo}}(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.Remove{{$txt.Function.NameGo}}(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.{{$txt.Function.NameGo}}(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	{{- if .ToJoinModel}}

	if len(b.R.{{$txt.Function.ForeignNameGo}}) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.{{$txt.Function.ForeignNameGo}}) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.{{$txt.Function.ForeignNameGo}}[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.{{$txt.Function.ForeignNameGo}}[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	{{- else}}

	if b.{{$txt.ForeignModel.ColumnNameGo}}.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.{{$txt.ForeignModel.ColumnNameGo}}.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.{{$txt.Function.ForeignNameGo}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.{{$txt.Function.ForeignNameGo}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.{{$txt.Function.ForeignNameGo}} != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.{{$txt.Function.ForeignNameGo}} != &a {
		t.Error("relationship to a should have been preserved")
	}
	{{- end}}

	if len(a.R.{{$txt.Function.NameGo}}) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a smodel deletion for performance so we have to flip the order
	if a.R.{{$txt.Function.NameGo}}[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.{{$txt.Function.NameGo}}[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}
{{end -}}
{{- end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* outer if join model */ -}}
