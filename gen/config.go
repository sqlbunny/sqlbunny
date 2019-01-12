package gen

import (
	"github.com/kernelpayments/sqlbunny/runtime/queries"
	"github.com/kernelpayments/sqlbunny/schema"
)

type ConfigItem interface {
	IsConfigItem()
}

type ConfigStruct struct {
	Items  []ConfigItem
	Schema *schema.Schema

	Dialect queries.Dialect

	PkgName    string
	OutputPath string
}

var Config *ConfigStruct
