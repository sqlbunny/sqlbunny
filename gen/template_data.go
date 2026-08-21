package gen

func BaseTemplateData() map[string]any {
	return map[string]any{
		"PkgName":     Config.ModelsPackageName,
		"Schema":      Config.Schema,
		"StringFuncs": templateStringMappers,
	}
}
