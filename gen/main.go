package gen

import (
	"log"
	"os"

	"github.com/kernelpayments/sqlbunny/def"
	"github.com/kernelpayments/sqlbunny/runtime/queries"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func Run(items ...def.ConfigItem) {
	schema, err := def.BuildSchema(items)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	rootCmd = &cobra.Command{Use: "sqlbunny"}

	Config = &ConfigStruct{
		Schema: schema,
		Items:  items,

		Dialect: queries.Dialect{
			LQ:                '"',
			RQ:                '"',
			IndexPlaceholders: true,
			UseTopClause:      false,
		},

		OutputPath: "models",
		PkgName:    "models",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use: "gen",
		Run: gen,
	})

	for _, i := range items {
		if p, ok := i.(Plugin); ok {
			p.BunnyPlugin()
		}
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func gen(cmd *cobra.Command, args []string) {
	for _, f := range genFuncs {
		f()
	}
}
