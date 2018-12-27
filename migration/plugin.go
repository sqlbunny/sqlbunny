package migration

import (
	"bytes"
	"fmt"
	"log"

	"github.com/kernelpayments/sqlbunny/config"
	"github.com/kernelpayments/sqlbunny/gen"
	"github.com/kernelpayments/sqlbunny/schema"
	"github.com/spf13/cobra"
)

type Plugin struct {
	Store *Store
}

func (Plugin) IsConfigItem() {}

func (p Plugin) InitPlugin() {
	config.Config.RootCmd.AddCommand(&cobra.Command{
		Use: "genmigrations",
		Run: p.cmdGenMigrations,
	})
}

func (p Plugin) RunPlugin() {
}

func (p Plugin) cmdGenMigrations(cmd *cobra.Command, args []string) {
	if p.Store == nil {
		log.Fatal("migrate.Plugin.Store is not set.")
	}

	s2 := config.Config.Schema

	s1 := schema.New()
	p.Store.applyAll(s1)

	ops := Diff(nil, s1, s2)

	if len(ops) == 0 {
		log.Fatal("No model changes found, doing nothing.")
	}
	migrationNumber := p.Store.nextFree()
	migrationFile := fmt.Sprintf("migration_%05d.go", migrationNumber)

	var buf bytes.Buffer
	gen.WritePackageName(&buf, "migrations")
	buf.WriteString("import \"github.com/kernelpayments/sqlbunny/migration\"\n")
	buf.WriteString(fmt.Sprintf("func init() {\nMigrations.Register(%d, ", migrationNumber))
	ops.Dump(&buf)
	buf.WriteString(")\n}")

	gen.WriteFile("migrations", migrationFile, &buf)
}
