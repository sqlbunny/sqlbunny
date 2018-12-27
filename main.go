package sqlbunny

import (
	"log"
	"os"

	"github.com/kernelpayments/sqlbunny/config"
	"github.com/kernelpayments/sqlbunny/def"
	"github.com/kernelpayments/sqlbunny/runtime/queries"
	"github.com/spf13/cobra"
)

type Plugin interface {
	InitPlugin()
	RunPlugin()
}

func Run(items ...def.ConfigItem) {
	schema, err := def.BuildSchema(items)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var rootCmd = &cobra.Command{Use: "sqlbunny"}

	config.Config = &config.ConfigStruct{
		Schema:  schema,
		Items:   items,
		RootCmd: rootCmd,

		Dialect: queries.Dialect{
			LQ:                '"',
			RQ:                '"',
			IndexPlaceholders: true,
			UseTopClause:      false,
		},

		OutputPath: "models",
		PkgName:    "models",
	}

	for _, i := range items {
		if p, ok := i.(Plugin); ok {
			p.InitPlugin()
		}
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
