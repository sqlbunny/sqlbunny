{{- if .Model.IsJoinModel -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Model.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $dot.Models $dot.Model . -}}
		{{- $varNameSingular := .Model | singular | camelCase -}}
		{{- $foreignVarNameSingular := .ForeignModel | singular | camelCase}}
func test{{$txt.LocalModel.NameGo}}OneToOne{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}}(t *testing.T) {
	tx := MustTx(bunny.Begin())
	defer tx.Rollback()

	var foreign {{$txt.ForeignModel.NameGo}}
	var local {{$txt.LocalModel.NameGo}}

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &foreign, {{$foreignVarNameSingular}}DBTypes, true, {{$foreignVarNameSingular}}FieldsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$txt.ForeignModel.NameGo}} struct: %s", err)
	}
	if err := randomize.Struct(seed, &local, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$txt.LocalModel.NameGo}} struct: %s", err)
	}

	{{if .ForeignColumnNullable -}}
	foreign.{{$txt.ForeignModel.ColumnNameGo}}.Valid = true
	{{- end}}
	{{if .Nullable -}}
	local.{{$txt.LocalModel.ColumnNameGo}}.Valid = true
	{{- end}}

	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	foreign.{{$txt.Function.ForeignAssignment}} = local.{{$txt.Function.LocalAssignment}}
	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.{{$txt.Function.NameGo}}(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	{{if $txt.Function.UsesBytes -}}
	if 0 != bytes.Compare(check.{{$txt.Function.ForeignAssignment}}, foreign.{{$txt.Function.ForeignAssignment}}) {
	{{else -}}
	if check.{{$txt.Function.ForeignAssignment}} != foreign.{{$txt.Function.ForeignAssignment}} {
	{{end -}}
		t.Errorf("want: %v, got %v", foreign.{{$txt.Function.ForeignAssignment}}, check.{{$txt.Function.ForeignAssignment}})
	}

	slice := {{$txt.LocalModel.NameGo}}Slice{&local}
	if err = local.L.Load{{$txt.Function.NameGo}}(tx, false, (*[]*{{$txt.LocalModel.NameGo}})(&slice)); err != nil {
		t.Fatal(err)
	}
	if local.R.{{$txt.Function.NameGo}} == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.{{$txt.Function.NameGo}} = nil
	if err = local.L.Load{{$txt.Function.NameGo}}(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.{{$txt.Function.NameGo}} == nil {
		t.Error("struct should have been eager loaded")
	}
}

{{end -}}{{/* range */}}
{{- end -}}{{/* join model */}}
