package migration

import (
	"bytes"
	"fmt"
	"log"

	"github.com/kernelpayments/sqlbunny/gen"
	"github.com/kernelpayments/sqlbunny/runtime/migration"
	"github.com/kernelpayments/sqlbunny/schema"
	"github.com/spf13/cobra"
)

type Plugin struct {
	Store *migration.Store
}

var _ gen.Plugin = &Plugin{}

func (*Plugin) IsConfigItem() {}

func (p *Plugin) BunnyPlugin() {
	gen.AddCommand(&cobra.Command{
		Use: "genmigrations",
		Run: p.cmdGenMigrations,
	})
}

func (p *Plugin) cmdGenMigrations(cmd *cobra.Command, args []string) {
	if p.Store == nil {
		log.Fatal("migrate.Plugin.Store is not set.")
	}

	s2 := gen.Config.Schema

	s1 := schema.New()
	p.applyAll(s1)

	ops := diff(nil, s1, s2)

	if len(ops) == 0 {
		log.Fatal("No model changes found, doing nothing.")
	}
	migrationNumber := p.nextFree()
	migrationFile := fmt.Sprintf("migration_%05d.go", migrationNumber)

	var buf bytes.Buffer
	gen.WritePackageName(&buf, "migrations")
	buf.WriteString("import \"github.com/kernelpayments/sqlbunny/migration\"\n")
	buf.WriteString(fmt.Sprintf("func init() {\nMigrations.Register(%d, ", migrationNumber))
	ops.Dump(&buf)
	buf.WriteString(")\n}")

	gen.WriteFile("migrations", migrationFile, buf.Bytes())
}

func (p *Plugin) applyAll(db *schema.Schema) {
	var i int64 = 1
	for {
		ops, ok := p.Store.Migrations[i]
		if !ok {
			break
		}

		for _, op := range ops {
			op.Apply(db)
		}

		i++
	}
}

func (p *Plugin) nextFree() int64 {
	var i int64 = 1
	for {
		_, ok := p.Store.Migrations[i]
		if !ok {
			break
		}

		i++
	}
	return i
}
