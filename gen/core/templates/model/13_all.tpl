{{- $modelName := .Model.Name | titleCase -}}
{{- $modelNamePlural := .Model.Name | plural | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase}}
// {{$modelNamePlural}} creates a {{$modelNamePlural}} query with the given mods.
func {{$modelNamePlural}}(mods ...qm.QueryMod) {{$varNameSingular}}Query {
	mods = append(mods, qm.From("{{.Model.Name | schemaModel}}"))
	return {{$varNameSingular}}Query{NewQuery(mods...)}
}
