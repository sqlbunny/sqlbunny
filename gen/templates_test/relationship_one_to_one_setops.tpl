{{- if .Model.IsJoinModel -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Model.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $dot.Models $dot.Model .}}
{{- $varNameSingular := .Model | singular | camelCase -}}
{{- $foreignVarNameSingular := .ForeignModel | singular | camelCase -}}
{{- $foreignPrimaryKeyCols := (getModel $dot.Models .ForeignModel).PrimaryKey.Columns}}
func test{{$txt.LocalModel.NameGo}}OneToOneSetOp{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}}(t *testing.T) {
	var err error

	tx := MustTx(bunny.Begin())
	defer tx.Rollback()

	var a {{$txt.LocalModel.NameGo}}
	var b, c {{$txt.ForeignModel.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, strmangle.SetComplement({{$varNameSingular}}PrimaryKeyColumns, {{$varNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, {{$foreignVarNameSingular}}DBTypes, false, strmangle.SetComplement({{$foreignVarNameSingular}}PrimaryKeyColumns, {{$foreignVarNameSingular}}FieldsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, {{$foreignVarNameSingular}}DBTypes, false, strmangle.SetComplement({{$foreignVarNameSingular}}PrimaryKeyColumns, {{$foreignVarNameSingular}}FieldsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*{{$txt.ForeignModel.NameGo}}{&b, &c} {
		err = a.Set{{$txt.Function.NameGo}}(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.{{$txt.Function.NameGo}} != x {
			t.Error("relationship struct not set to correct value")
		}
		if x.R.{{$txt.Function.ForeignNameGo}} != &a {
			t.Error("failed to append to foreign relationship struct")
		}

		{{if $txt.Function.UsesBytes -}}
		if 0 != bytes.Compare(a.{{$txt.Function.LocalAssignment}}, x.{{$txt.Function.ForeignAssignment}}) {
		{{else -}}
		if a.{{$txt.Function.LocalAssignment}} != x.{{$txt.Function.ForeignAssignment}} {
		{{end -}}
			t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}})
		}

		{{if setInclude .ForeignColumn $foreignPrimaryKeyCols -}}
		if exists, err := {{$txt.ForeignModel.NameGo}}Exists(tx, x.{{$foreignPrimaryKeyCols | stringMap $dot.StringFuncs.titleCase | join ", x."}}); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'x' to exist")
		}
		{{else -}}
		zero := reflect.Zero(reflect.TypeOf(x.{{$txt.Function.ForeignAssignment}}))
		reflect.Indirect(reflect.ValueOf(&x.{{$txt.Function.ForeignAssignment}})).Set(zero)

		if err = x.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}
		{{- end}}

		{{if $txt.Function.UsesBytes -}}
		if 0 != bytes.Compare(a.{{$txt.Function.LocalAssignment}}, x.{{$txt.Function.ForeignAssignment}}) {
		{{else -}}
		if a.{{$txt.Function.LocalAssignment}} != x.{{$txt.Function.ForeignAssignment}} {
		{{end -}}
			t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, x.{{$txt.Function.ForeignAssignment}})
		}

		if err = x.Delete(tx); err != nil {
			t.Fatal("failed to delete x", err)
		}
	}
}
{{- if .ForeignColumnNullable}}

func test{{$txt.LocalModel.NameGo}}OneToOneRemoveOp{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}}(t *testing.T) {
	var err error

	tx := MustTx(bunny.Begin())
	defer tx.Rollback()

	var a {{$txt.LocalModel.NameGo}}
	var b {{$txt.ForeignModel.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, strmangle.SetComplement({{$varNameSingular}}PrimaryKeyColumns, {{$varNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, {{$foreignVarNameSingular}}DBTypes, false, strmangle.SetComplement({{$foreignVarNameSingular}}PrimaryKeyColumns, {{$foreignVarNameSingular}}FieldsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.Set{{$txt.Function.NameGo}}(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.Remove{{$txt.Function.NameGo}}(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.{{$txt.Function.NameGo}}(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.{{$txt.Function.NameGo}} != nil {
		t.Error("R struct entry should be nil")
	}

	if b.{{$txt.ForeignModel.ColumnNameGo}}.Valid {
		t.Error("foreign key field should be nil")
	}

	if b.R.{{$txt.Function.ForeignNameGo}} != nil {
		t.Error("failed to remove a from b's relationships")
	}
}
{{end -}}{{/* end if foreign key nullable */}}
{{- end -}}{{/* range */}}
{{- end -}}{{/* join model */}}
