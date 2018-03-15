// Package main defines a command line interface for the sqlboiler package
package main

import (
	"fmt"
	"os"

	"github.com/KernelPay/sqlboiler/gen"
	"github.com/KernelPay/sqlboiler/migration"
	"github.com/KernelPay/sqlboiler/schema"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const sqlBoilerVersion = "2.6.0"

var cmdConfig gen.Config
var schemaFile string

func main() {
	var cmdGen = &cobra.Command{
		Use:  "gen [flags]",
		RunE: runGen,
	}
	var cmdMigrationGen = &cobra.Command{
		Use:  "gen",
		RunE: runMigrationGen,
	}
	cmdMigrationGen.Flags().StringVarP(&schemaFile, "schema", "i", "", "Schema file to load")
	var cmdMigrationRun = &cobra.Command{
		Use:  "run",
		RunE: runMigrationRun,
	}

	var cmdMigration = &cobra.Command{
		Use: "migration",
	}
	cmdMigration.AddCommand(cmdMigrationGen, cmdMigrationRun)

	var cmdVersion = &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("SQLBoiler v" + sqlBoilerVersion)
		},
	}

	var rootCmd = &cobra.Command{Use: "sqlboiler"}
	rootCmd.AddCommand(cmdGen, cmdVersion, cmdMigration)

	// Set up the cobra root command flags
	cmdGen.Flags().StringVarP(&schemaFile, "schema", "i", "", "Schema file to load")
	cmdGen.Flags().StringVarP(&cmdConfig.OutFolder, "output", "o", "models", "The name of the folder to output to")
	cmdGen.Flags().StringVarP(&cmdConfig.PkgName, "pkgname", "p", "models", "The name you wish to assign to your generated package")
	cmdGen.Flags().StringVarP(&cmdConfig.BaseDir, "basedir", "", "", "The base directory has the templates and templates_test folders")
	cmdGen.Flags().StringSliceVarP(&cmdConfig.Tags, "tags", "t", nil, "Struct tags to be included on your models in addition to json, yaml, toml")
	cmdGen.Flags().BoolVarP(&cmdConfig.NoTests, "no-tests", "", false, "Disable generated go test files")
	cmdGen.Flags().BoolVarP(&cmdConfig.NoHooks, "no-hooks", "", false, "Disable hooks feature for your models")
	cmdGen.Flags().BoolVarP(&cmdConfig.Wipe, "wipe", "", false, "Delete the output folder (rm -rf) before generation to ensure sanity")

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

// initModels retrieves all "public" schema model names from the database.
func loadSchema(filename string) (*schema.Schema, error) {
	var err error
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open schema file")
	}
	schema, err := schema.ParseSchema(f)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse schema")
	}

	return schema, nil
}

func runGen(cmd *cobra.Command, args []string) error {
	schema, err := loadSchema(schemaFile)
	if err != nil {
		return err
	}

	cmdState, err := gen.New(schema, &cmdConfig)
	if err != nil {
		return err
	}

	return cmdState.Run(true)
}

func runMigrationGen(cmd *cobra.Command, args []string) error {
	schema, err := loadSchema(schemaFile)
	if err != nil {
		return err
	}

	return migration.Run(schema)
}

func runMigrationRun(cmd *cobra.Command, args []string) error {
	return nil
}
