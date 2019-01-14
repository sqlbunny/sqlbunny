package gen

import (
	"path/filepath"

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

	OutputPath        string
	ModelsPackageName string
}

func (c *ConfigStruct) ModelsOutputPath() string {
	return filepath.Join(c.OutputPath, c.ModelsPackageName)
}

var Config *ConfigStruct
