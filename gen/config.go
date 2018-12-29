package gen

import (
	"github.com/kernelpayments/sqlbunny/def"
	"github.com/kernelpayments/sqlbunny/runtime/queries"
	"github.com/kernelpayments/sqlbunny/schema"
)

type ConfigStruct struct {
	Items  []def.ConfigItem
	Schema *schema.Schema

	Dialect queries.Dialect

	PkgName    string
	OutputPath string
}

var Config *ConfigStruct
