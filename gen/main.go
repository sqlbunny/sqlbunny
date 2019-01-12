package gen

import (
	"os"

	"github.com/kernelpayments/sqlbunny/runtime/queries"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

type expander interface {
	Expand() []ConfigItem
}

func expand(items []ConfigItem, item ConfigItem) []ConfigItem {
	items = append(items, item)
	if e, ok := item.(expander); ok {
		items = append(items, e.Expand()...)
	}
	return items
}

func expandAll(items []ConfigItem) []ConfigItem {
	var res []ConfigItem
	for _, i := range items {
		res = expand(res, i)
	}
	return res
}

func Run(items []ConfigItem) {
	items = expandAll(items)

	rootCmd = &cobra.Command{Use: "sqlbunny"}

	Config = &ConfigStruct{
		Items: items,

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
