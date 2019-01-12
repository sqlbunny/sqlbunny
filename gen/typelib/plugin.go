package typelib

import (
	"github.com/kernelpayments/sqlbunny/gen"
	"github.com/kernelpayments/sqlbunny/gen/core"
)

type Plugin struct {
}

var _ gen.Plugin = &Plugin{}

func (*Plugin) IsConfigItem() {}

func (p *Plugin) BunnyPlugin() {
}

func (p *Plugin) Expand() []gen.ConfigItem {
	return []gen.ConfigItem{
		core.Type("int32", core.BaseType{
			Go:       "int32",
			GoNull:   "github.com/kernelpayments/sqlbunny/types/null.Int32",
			Postgres: "integer",
		}),

		core.Type("int64", core.BaseType{
			Go:       "int64",
			GoNull:   "github.com/kernelpayments/sqlbunny/types/null.Int64",
			Postgres: "bigint",
		}),

		core.Type("float32", core.BaseType{
			Go:       "float32",
			GoNull:   "github.com/kernelpayments/sqlbunny/types/null.Float32",
			Postgres: "real",
		}),

		core.Type("float64", core.BaseType{
			Go:       "float64",
			GoNull:   "github.com/kernelpayments/sqlbunny/types/null.Float64",
			Postgres: "double precision",
		}),

		core.Type("bool", core.BaseType{
			Go:       "bool",
			GoNull:   "github.com/kernelpayments/sqlbunny/types/null.Bool",
			Postgres: "boolean",
		}),

		core.Type("string", core.BaseType{
			Go:       "string",
			GoNull:   "github.com/kernelpayments/sqlbunny/types/null.String",
			Postgres: "text",
		}),

		core.Type("bytea", core.BaseType{
			Go:       "[]byte",
			GoNull:   "github.com/kernelpayments/sqlbunny/types/null.Bytes",
			Postgres: "bytea",
		}),

		core.Type("jsonb", core.BaseType{
			Go:       "github.com/kernelpayments/sqlbunny/types.JSON",
			GoNull:   "github.com/kernelpayments/sqlbunny/types/null.JSON",
			Postgres: "jsonb",
		}),

		core.Type("time", core.BaseType{
			Go:       "time.Time",
			GoNull:   "github.com/kernelpayments/sqlbunny/types/null.Time",
			Postgres: "timestamptz",
		}),
	}
}
