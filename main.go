package sqlbunny

import (
	"fmt"
	"os"

	"github.com/kernelpayments/sqlbunny/def"
	"github.com/kernelpayments/sqlbunny/gen"
	"github.com/kernelpayments/sqlbunny/migration"
	"github.com/spf13/cobra"
)

const sqlbunnyVersion = "2.6.0"

var cmdConfig gen.Config

func Run() {
	var cmdGen = &cobra.Command{
		Use:  "gen [flags]",
		RunE: runGen,
	}
	cmdGen.Flags().StringVarP(&cmdConfig.OutFolder, "output", "o", "models", "The name of the folder to output to")
	cmdGen.Flags().StringVarP(&cmdConfig.PkgName, "pkgname", "p", "models", "The name you wish to assign to your generated package")
	cmdGen.Flags().StringSliceVarP(&cmdConfig.Tags, "tags", "t", nil, "Struct tags to be included on your models in addition to json, yaml, toml")
	cmdGen.Flags().BoolVarP(&cmdConfig.NoTests, "no-tests", "", false, "Disable generated go test files")
	cmdGen.Flags().BoolVarP(&cmdConfig.NoHooks, "no-hooks", "", false, "Disable hooks feature for your models")
	cmdGen.Flags().BoolVarP(&cmdConfig.Wipe, "wipe", "", false, "Delete the output folder (rm -rf) before generation to ensure sanity")

	var cmdGenMigrations = &cobra.Command{
		Use:  "genmigrations",
		RunE: runGenMigrations,
	}
	var cmdVersion = &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("sqlbunny v" + sqlbunnyVersion)
		},
	}

	var rootCmd = &cobra.Command{Use: "sqlbunny"}
	rootCmd.AddCommand(cmdGen, cmdGenMigrations, cmdVersion)

	// hide flags not recommended for use
	rootCmd.PersistentFlags().MarkHidden("replace")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type commandFailure string

func (c commandFailure) Error() string {
	return string(c)
}

func runGen(cmd *cobra.Command, args []string) error {
	schema, err := def.Schema()
	if err != nil {
		return err
	}

	cmdState, err := gen.New(schema, &cmdConfig)
	if err != nil {
		return err
	}

	return cmdState.Run(true)
}

var mstore *migration.MigrationStore

func SetMigrations(s *migration.MigrationStore) {
	mstore = s
}

func runGenMigrations(cmd *cobra.Command, args []string) error {
	schema, err := def.Schema()
	if err != nil {
		return err
	}

	return migration.Run(schema, mstore)
}
