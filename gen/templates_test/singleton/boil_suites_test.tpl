{{- $dot := .}}
// This test suite runs each operation test in parallel.
// Example, if your database has 3 models, the suite will run:
// model1, model2 and model3 Delete in parallel
// model1, model2 and model3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}})
  {{end -}}
  {{- end -}}
}

func TestDelete(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}Delete)
  {{end -}}
  {{- end -}}
}

func TestQueryDeleteAll(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}QueryDeleteAll)
  {{end -}}
  {{- end -}}
}

func TestSliceDeleteAll(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}SliceDeleteAll)
  {{end -}}
  {{- end -}}
}

func TestExists(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}Exists)
  {{end -}}
  {{- end -}}
}

func TestFind(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}Find)
  {{end -}}
  {{- end -}}
}

func TestBind(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}Bind)
  {{end -}}
  {{- end -}}
}

func TestOne(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}One)
  {{end -}}
  {{- end -}}
}

func TestAll(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}All)
  {{end -}}
  {{- end -}}
}

func TestCount(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}Count)
  {{end -}}
  {{- end -}}
}

{{if not .NoHooks -}}
func TestHooks(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}Hooks)
  {{end -}}
  {{- end -}}
}
{{- end}}

func TestInsert(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}Insert)
  t.Run("{{$modelName}}", test{{$modelName}}InsertWhitelist)
  {{end -}}
  {{- end -}}
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
{{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
    {{- range $model.ForeignKeys -}}
      {{- $txt := txtsFromFKey $dot.Models $model . -}}
  t.Run("{{$txt.LocalModel.NameGo}}To{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}}", test{{$txt.LocalModel.NameGo}}ToOne{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}})
    {{end -}}{{- /* fkey range */ -}}
  {{- end -}}{{- /* if join model */ -}}
{{- end -}}{{- /* models range */ -}}
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {
  {{- range $index, $model := .Models}}
	{{- if $model.IsJoinModel -}}
	{{- else -}}
	  {{- range $model.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $dot.Models $model . -}}
  t.Run("{{$txt.LocalModel.NameGo}}To{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}}", test{{$txt.LocalModel.NameGo}}OneToOne{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}})
	  {{end -}}{{- /* range */ -}}
	{{- end -}}{{- /* outer if join model */ -}}
  {{- end -}}{{- /* outer models range */ -}}
}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
  {{- range $index, $model := .Models}}
    {{- if $model.IsJoinModel -}}
    {{- else -}}
      {{- range $model.ToManyRelationships -}}
        {{- $txt := txtsFromToMany $dot.Models $model . -}}
  t.Run("{{$txt.LocalModel.NameGo}}To{{$txt.Function.NameGo}}", test{{$txt.LocalModel.NameGo}}ToMany{{$txt.Function.NameGo}})
      {{end -}}{{- /* range */ -}}
    {{- end -}}{{- /* outer if join model */ -}}
  {{- end -}}{{- /* outer models range */ -}}
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
{{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
    {{- range $model.ForeignKeys -}}
      {{- $txt := txtsFromFKey $dot.Models $model . -}}
  t.Run("{{$txt.LocalModel.NameGo}}To{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}}", test{{$txt.LocalModel.NameGo}}ToOneSetOp{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}})
    {{end -}}{{- /* fkey range */ -}}
  {{- end -}}{{- /* if join model */ -}}
{{- end -}}{{- /* models range */ -}}
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
{{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
    {{- range $model.ForeignKeys -}}
      {{- $txt := txtsFromFKey $dot.Models $model . -}}
      {{- if $txt.ForeignKey.Nullable -}}
  t.Run("{{$txt.LocalModel.NameGo}}To{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}}", test{{$txt.LocalModel.NameGo}}ToOneRemoveOp{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}})
      {{end -}}{{- /* if foreign key nullable */ -}}
    {{- end -}}{{- /* fkey range */ -}}
  {{- end -}}{{- /* if join model */ -}}
{{- end -}}{{- /* models range */ -}}
}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {
  {{- range $index, $model := .Models}}
	{{- if $model.IsJoinModel -}}
	{{- else -}}
	  {{- range $model.ToOneRelationships -}}
		  {{- $txt := txtsFromOneToOne $dot.Models $model . -}}
	t.Run("{{$txt.LocalModel.NameGo}}To{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}}", test{{$txt.LocalModel.NameGo}}OneToOneSetOp{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}})
	  {{end -}}{{- /* range to one relationships */ -}}
	{{- end -}}{{- /* outer if join model */ -}}
  {{- end -}}{{- /* outer models range */ -}}
}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {
  {{- range $index, $model := .Models}}
	{{- if $model.IsJoinModel -}}
	{{- else -}}
	  {{- range $model.ToOneRelationships -}}
		{{- if .ForeignColumnNullable -}}
		  {{- $txt := txtsFromOneToOne $dot.Models $model . -}}
	t.Run("{{$txt.LocalModel.NameGo}}To{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}}", test{{$txt.LocalModel.NameGo}}OneToOneRemoveOp{{$txt.ForeignModel.NameGo}}Using{{$txt.Function.NameGo}})
		{{end -}}{{- /* if foreign field nullable */ -}}
	  {{- end -}}{{- /* range */ -}}
	{{- end -}}{{- /* outer if join model */ -}}
  {{- end -}}{{- /* outer models range */ -}}
}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
  {{- range $index, $model := .Models}}
    {{- if $model.IsJoinModel -}}
    {{- else -}}
      {{- range $model.ToManyRelationships -}}
        {{- $txt := txtsFromToMany $dot.Models $model . -}}
  t.Run("{{$txt.LocalModel.NameGo}}To{{$txt.Function.NameGo}}", test{{$txt.LocalModel.NameGo}}ToManyAddOp{{$txt.Function.NameGo}})
      {{end -}}{{- /* range */ -}}
    {{- end -}}{{- /* outer if join model */ -}}
  {{- end -}}{{- /* outer models range */ -}}
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
  {{- range $index, $model := .Models}}
    {{- if $model.IsJoinModel -}}
    {{- else -}}
      {{- range $model.ToManyRelationships -}}
        {{- if not (or .ForeignColumnNullable .ToJoinModel)}}
        {{- else -}}
          {{- $txt := txtsFromToMany $dot.Models $model . -}}
    t.Run("{{$txt.LocalModel.NameGo}}To{{$txt.Function.NameGo}}", test{{$txt.LocalModel.NameGo}}ToManySetOp{{$txt.Function.NameGo}})
        {{end -}}{{- /* if foreign field nullable */ -}}
      {{- end -}}{{- /* range */ -}}
    {{- end -}}{{- /* outer if join model */ -}}
  {{- end -}}{{- /* outer models range */ -}}
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
  {{- range $index, $model := .Models}}
    {{- if $model.IsJoinModel -}}
    {{- else -}}
      {{- range $model.ToManyRelationships -}}
        {{- if not (or .ForeignColumnNullable .ToJoinModel)}}
        {{- else -}}
          {{- $txt := txtsFromToMany $dot.Models $model . -}}
    t.Run("{{$txt.LocalModel.NameGo}}To{{$txt.Function.NameGo}}", test{{$txt.LocalModel.NameGo}}ToManyRemoveOp{{$txt.Function.NameGo}})
        {{end -}}{{- /* if foreign field nullable */ -}}
      {{- end -}}{{- /* range */ -}}
    {{- end -}}{{- /* outer if join model */ -}}
  {{- end -}}{{- /* outer models range */ -}}
}

func TestReload(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}Reload)
  {{end -}}
  {{- end -}}
}

func TestReloadAll(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}ReloadAll)
  {{end -}}
  {{- end -}}
}

func TestSelect(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}Select)
  {{end -}}
  {{- end -}}
}

func TestUpdate(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}Update)
  {{end -}}
  {{- end -}}
}

func TestSliceUpdateAll(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}SliceUpdateAll)
  {{end -}}
  {{- end -}}
}

func TestUpsert(t *testing.T) {
  {{- range $index, $model := .Models}}
  {{- if $model.IsJoinModel -}}
  {{- else -}}
  {{- $modelName := $model.Name | plural | titleCase -}}
  t.Run("{{$modelName}}", test{{$modelName}}Upsert)
  {{end -}}
  {{- end -}}
}
