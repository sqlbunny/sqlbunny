package migration

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/kernelpayments/sqlbunny/common"
	"github.com/kernelpayments/sqlbunny/def"
	"github.com/kernelpayments/sqlbunny/schema"
	"github.com/spf13/cobra"
)

type Plugin struct {
	Store *Store
}

func (Plugin) IsConfigItem() {}

func (p Plugin) AddCommands(config *def.Config, rootCmd *cobra.Command) {
	var cmd = &cobra.Command{
		Use: "genmigrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			if p.Store == nil {
				return errors.New("migrate.Plugin.Store is not set.")
			}

			s2 := config.Schema

			s1 := schema.New()
			p.Store.applyAll(s1)

			ops := Diff(nil, s1, s2)

			if len(ops) == 0 {
				return fmt.Errorf("No model changes found, doing nothing.")
			}
			migrationNumber := p.Store.nextFree()
			migrationFile := fmt.Sprintf("migration_%05d.go", migrationNumber)

			var buf bytes.Buffer
			common.WritePackageName(&buf, "migrations")
			buf.WriteString("import \"github.com/kernelpayments/sqlbunny/migration\"\n")
			buf.WriteString(fmt.Sprintf("func init() {\nMigrations.Register(%d, ", migrationNumber))
			ops.Dump(&buf)
			buf.WriteString(")\n}")

			if err := common.WriteFile("migrations", migrationFile, &buf); err != nil {
				return err
			}
			return nil
		},
	}
	rootCmd.AddCommand(cmd)
}
