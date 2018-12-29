package gen

import (
	"github.com/kernelpayments/sqlbunny/def"
	"github.com/kernelpayments/sqlbunny/runtime/queries"
	"github.com/kernelpayments/sqlbunny/schema"
	"github.com/spf13/cobra"
)

type ConfigStruct struct {
	Items   []def.ConfigItem
	Schema  *schema.Schema
	RootCmd *cobra.Command

	Dialect queries.Dialect

	PkgName    string
	OutputPath string
}

var Config *ConfigStruct
