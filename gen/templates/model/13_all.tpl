{{- $modelNamePlural := .Model.Name | plural | titleCase -}}
{{- $varNameSingular := .Model.Name | singular | camelCase}}
// {{$modelNamePlural}} retrieves all the records using an executor.
func {{$modelNamePlural}}(ctx context.Context, mods ...qm.QueryMod) {{$varNameSingular}}Query {
	mods = append(mods, qm.From("{{.Model.Name | .SchemaModel}}"))
	return {{$varNameSingular}}Query{NewQuery(ctx, mods...)}
}
