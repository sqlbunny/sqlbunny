{{- if .Model.IsJoinModel -}}
{{- else -}}
	{{- $dot := . }}
	{{- $model := .Model }}
	{{- range .Model.ToManyRelationships -}}
	{{- $txt := txtsFromToMany $dot.Models $model .}}
	{{- $varNameSingular := .Model | singular | camelCase -}}
	{{- $foreignVarNameSingular := .ForeignModel | singular | camelCase -}}
func test{{$txt.LocalModel.NameGo}}ToMany{{$txt.Function.NameGo}}(t *testing.T) {
	var err error
	tx := MustTx(bunny.Begin())
	defer tx.Rollback()

	var a {{$txt.LocalModel.NameGo}}
	var b, c {{$txt.ForeignModel.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$txt.LocalModel.NameGo}} struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, {{$foreignVarNameSingular}}DBTypes, false, {{$foreignVarNameSingular}}FieldsWithDefault...)
	randomize.Struct(seed, &c, {{$foreignVarNameSingular}}DBTypes, false, {{$foreignVarNameSingular}}FieldsWithDefault...)
	{{if .Nullable -}}
	a.{{.Field | titleCase}}.Valid = true
	{{- end}}
	{{- if .ForeignColumnNullable}}
	b.{{.ForeignColumn | titleCase}}.Valid = true
	c.{{.ForeignColumn | titleCase}}.Valid = true
	{{- end}}
	{{if not .ToJoinModel -}}
	b.{{$txt.Function.ForeignAssignment}} = a.{{$txt.Function.LocalAssignment}}
	c.{{$txt.Function.ForeignAssignment}} = a.{{$txt.Function.LocalAssignment}}
	{{- end}}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	{{if .ToJoinModel -}}
	_, err = tx.Exec("insert into {{.JoinModel | $dot.SchemaModel}} ({{.JoinLocalColumn | $dot.Quotes}}, {{.JoinForeignColumn | $dot.Quotes}}) values {{if $dot.Dialect.IndexPlaceholders}}($1, $2){{else}}(?, ?){{end}}", a.{{$txt.LocalModel.ColumnNameGo}}, b.{{$txt.ForeignModel.ColumnNameGo}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into {{.JoinModel | $dot.SchemaModel}} ({{.JoinLocalColumn | $dot.Quotes}}, {{.JoinForeignColumn | $dot.Quotes}}) values {{if $dot.Dialect.IndexPlaceholders}}($1, $2){{else}}(?, ?){{end}}", a.{{$txt.LocalModel.ColumnNameGo}}, c.{{$txt.ForeignModel.ColumnNameGo}})
	if err != nil {
		t.Fatal(err)
	}
	{{end}}

	{{$varname := .ForeignModel | singular | camelCase -}}
	{{$varname}}, err := a.{{$txt.Function.NameGo}}(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range {{$varname}} {
		{{if $txt.Function.UsesBytes -}}
		if 0 == bytes.Compare(v.{{$txt.Function.ForeignAssignment}}, b.{{$txt.Function.ForeignAssignment}}) {
			bFound = true
		}
		if 0 == bytes.Compare(v.{{$txt.Function.ForeignAssignment}}, c.{{$txt.Function.ForeignAssignment}}) {
			cFound = true
		}
		{{else -}}
		if v.{{$txt.Function.ForeignAssignment}} == b.{{$txt.Function.ForeignAssignment}} {
			bFound = true
		}
		if v.{{$txt.Function.ForeignAssignment}} == c.{{$txt.Function.ForeignAssignment}} {
			cFound = true
		}
		{{end -}}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := {{$txt.LocalModel.NameGo}}Slice{&a}
	if err = a.L.Load{{$txt.Function.NameGo}}(tx, false, (*[]*{{$txt.LocalModel.NameGo}})(&slice)); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.{{$txt.Function.NameGo}}); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.{{$txt.Function.NameGo}} = nil
	if err = a.L.Load{{$txt.Function.NameGo}}(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.{{$txt.Function.NameGo}}); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", {{$varname}})
	}
}

{{end -}}{{- /* range */ -}}
{{- end -}}{{- /* outer if join model */ -}}
