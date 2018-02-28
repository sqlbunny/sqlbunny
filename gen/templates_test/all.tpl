{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
{{- $modelNamePlural := .Model.Name | plural | titleCase -}}
{{- $varNamePlural := .Model.Name | plural | camelCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
func test{{$modelNamePlural}}(t *testing.T) {
	t.Parallel()

	query := {{$modelNamePlural}}(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
