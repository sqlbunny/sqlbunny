package gen

import (
	"github.com/sqlbunny/sqlbunny/runtime/strmangle"
)

func BaseTemplateData() map[string]any {
	d := Config.Dialect
	lq := strmangle.QuoteCharacter(d.LQ)
	rq := strmangle.QuoteCharacter(d.RQ)

	return map[string]any{
		"PkgName":     Config.ModelsPackageName,
		"Schema":      Config.Schema,
		"Dialect":     d,
		"LQ":          lq,
		"RQ":          rq,
		"StringFuncs": templateStringMappers,
	}
}
