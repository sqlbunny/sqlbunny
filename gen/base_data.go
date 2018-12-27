package gen

import (
	"fmt"

	"github.com/kernelpayments/sqlbunny/config"
	"github.com/kernelpayments/sqlbunny/runtime/strmangle"

	"github.com/kernelpayments/sqlbunny/runtime/queries"
)

// templateData for sqlbunny templates
type TemplateData struct {
	// Controls what names are output
	PkgName string
	Schema  string

	// Controls which code is output (mysql vs postgres ...)
	UseLastInsertID bool

	// Turn off hook generation
	NoHooks bool

	// StringFuncs are usable in templates with stringMap
	StringFuncs map[string]func(string) string

	// Dialect controls quoting
	Dialect queries.Dialect
	LQ      string
	RQ      string
}

func (t TemplateData) Quotes(s string) string {
	return fmt.Sprintf("%s%s%s", t.LQ, s, t.RQ)
}

func (t TemplateData) SchemaModel(model string) string {
	return strmangle.SchemaModel(t.LQ, t.RQ, model)
}

func BaseTemplateData() *TemplateData {
	return &TemplateData{
		UseLastInsertID: true,
		PkgName:         config.Config.PkgName,
		NoHooks:         false,
		Dialect:         config.Config.Dialect,
		LQ:              strmangle.QuoteCharacter(config.Config.Dialect.LQ),
		RQ:              strmangle.QuoteCharacter(config.Config.Dialect.RQ),

		StringFuncs: templateStringMappers,
	}
}
