package sqlbunny

import (
	"log"
	"os"

	"github.com/kernelpayments/sqlbunny/def"
	"github.com/kernelpayments/sqlbunny/gen"
	"github.com/spf13/cobra"
)

type commandPlugin interface {
	AddCommands(config *def.Config, rootCmd *cobra.Command)
}

func Run(items ...def.ConfigItem) {
	config, err := def.Validate(items)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var rootCmd = &cobra.Command{Use: "sqlbunny"}

	addGenCommand(config, rootCmd)

	for _, i := range items {
		if p, ok := i.(commandPlugin); ok {
			p.AddCommands(config, rootCmd)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func addGenCommand(config *def.Config, rootCmd *cobra.Command) {
	var cmdConfig gen.Config
	var cmd = &cobra.Command{
		Use: "gen [flags]",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdState, err := gen.New(config.Schema, &cmdConfig)
			if err != nil {
				return err
			}

			return cmdState.Run(true)
		},
	}
	cmd.Flags().StringVarP(&cmdConfig.OutFolder, "output", "o", "models", "The name of the folder to output to")
	cmd.Flags().StringVarP(&cmdConfig.PkgName, "pkgname", "p", "models", "The name you wish to assign to your generated package")
	cmd.Flags().StringSliceVarP(&cmdConfig.Tags, "tags", "t", nil, "Struct tags to be included on your models in addition to json, yaml, toml")
	cmd.Flags().BoolVarP(&cmdConfig.NoTests, "no-tests", "", false, "Disable generated go test files")
	cmd.Flags().BoolVarP(&cmdConfig.NoHooks, "no-hooks", "", false, "Disable hooks feature for your models")
	cmd.Flags().BoolVarP(&cmdConfig.Wipe, "wipe", "", false, "Delete the output folder (rm -rf) before generation to ensure sanity")
	rootCmd.AddCommand(cmd)
}
