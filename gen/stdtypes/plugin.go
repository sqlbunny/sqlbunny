package stdtypes

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
		core.Type("int16", core.BaseType{
			Go:     "int16",
			GoNull: "github.com/kernelpayments/sqlbunny/types/null.Int16",
			Postgres: core.SQLType{
				Type:      "smallint",
				ZeroValue: "0",
			},
		}),

		core.Type("int32", core.BaseType{
			Go:     "int32",
			GoNull: "github.com/kernelpayments/sqlbunny/types/null.Int32",
			Postgres: core.SQLType{
				Type:      "integer",
				ZeroValue: "0",
			},
		}),

		core.Type("int64", core.BaseType{
			Go:     "int64",
			GoNull: "github.com/kernelpayments/sqlbunny/types/null.Int64",
			Postgres: core.SQLType{
				Type:      "bigint",
				ZeroValue: "0",
			},
		}),

		core.Type("float32", core.BaseType{
			Go:     "float32",
			GoNull: "github.com/kernelpayments/sqlbunny/types/null.Float32",
			Postgres: core.SQLType{
				Type:      "real",
				ZeroValue: "0",
			},
		}),

		core.Type("float64", core.BaseType{
			Go:     "float64",
			GoNull: "github.com/kernelpayments/sqlbunny/types/null.Float64",
			Postgres: core.SQLType{
				Type:      "double precision",
				ZeroValue: "0",
			},
		}),

		core.Type("bool", core.BaseType{
			Go:     "bool",
			GoNull: "github.com/kernelpayments/sqlbunny/types/null.Bool",
			Postgres: core.SQLType{
				Type:      "boolean",
				ZeroValue: "false",
			},
		}),

		core.Type("string", core.BaseType{
			Go:     "string",
			GoNull: "github.com/kernelpayments/sqlbunny/types/null.String",
			Postgres: core.SQLType{
				Type:      "text",
				ZeroValue: "''",
			},
		}),

		core.Type("bytea", core.BaseType{
			Go:     "[]byte",
			GoNull: "github.com/kernelpayments/sqlbunny/types/null.Bytes",
			Postgres: core.SQLType{
				Type:      "bytea",
				ZeroValue: "''",
			},
		}),

		core.Type("jsonb", core.BaseType{
			Go:     "github.com/kernelpayments/sqlbunny/types.JSON",
			GoNull: "github.com/kernelpayments/sqlbunny/types/null.JSON",
			Postgres: core.SQLType{
				Type:      "jsonb",
				ZeroValue: "'null'",
			},
		}),

		core.Type("time", core.BaseType{
			Go:     "time.Time",
			GoNull: "github.com/kernelpayments/sqlbunny/types/null.Time",
			Postgres: core.SQLType{
				Type:      "timestamptz",
				ZeroValue: "",
			},
		}),
	}
}
