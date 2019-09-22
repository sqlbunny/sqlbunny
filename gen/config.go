package gen

import (
	"github.com/sqlbunny/sqlbunny/runtime/queries"
	"github.com/sqlbunny/sqlbunny/schema"
)

type ConfigItem interface {
	ConfigItem(ctx *Context)
}

type ConfigStruct struct {
	Items  []ConfigItem
	Schema *schema.Schema

	Dialect queries.Dialect

	ModelsPackagePath string
	ModelsPackageName string
}

var Config *ConfigStruct
