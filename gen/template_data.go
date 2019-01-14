package gen

import (
	"github.com/kernelpayments/sqlbunny/runtime/strmangle"
)

func BaseTemplateData() map[string]interface{} {
	d := Config.Dialect
	lq := strmangle.QuoteCharacter(d.LQ)
	rq := strmangle.QuoteCharacter(d.RQ)

	return map[string]interface{}{
		"PkgName":     Config.ModelsPackageName,
		"Dialect":     d,
		"LQ":          lq,
		"RQ":          rq,
		"StringFuncs": templateStringMappers,
	}
}
