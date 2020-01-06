package migration

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
	"github.com/sqlbunny/sqlbunny/gen"
	"github.com/sqlbunny/sqlbunny/runtime/migration"
	"github.com/sqlbunny/sqlschema/diff"
	"github.com/sqlbunny/sqlschema/schema"
)

type Plugin struct {
	Store       *migration.Store
	PackageName string
	PackagePath string
}

var _ gen.Plugin = &Plugin{}

func (*Plugin) ConfigItem(ctx *gen.Context) {}

func (p *Plugin) BunnyPlugin() {
	if p.PackageName == "" {
		p.PackageName = "migrations"
	}
	if p.PackagePath == "" {
		p.PackagePath = "./migrations"
	}

	cmd := &cobra.Command{
		Use: "migration",
	}
	gen.AddCommand(cmd)

	cmd.AddCommand(&cobra.Command{
		Use: "gen",
		Run: p.cmdGen,
	})
	cmd.AddCommand(&cobra.Command{
		Use: "check",
		Run: p.cmdCheck,
	})
	cmd.AddCommand(&cobra.Command{
		Use: "merge",
		Run: p.cmdMerge,
	})
	cmd.AddCommand(&cobra.Command{
		Use: "gensql",
		Run: p.cmdGenSQL,
	})
}

func (p *Plugin) cmdCheck(cmd *cobra.Command, args []string) {
	if p.Store == nil {
		log.Fatal("migrate.Plugin.Store is not set.")
	}

	s1 := schema.NewDatabase()
	p.applyAll(s1)
	s2 := gen.Config.Schema.SQLSchema()
	ops := diff.Diff(s1, s2)

	if len(ops) != 0 {
		log.Fatal("Migrations are not up to date with the defined models. You need to run 'migration gen'.")
	}
}

func (p *Plugin) cmdMerge(cmd *cobra.Command, args []string) {
	p.ensureStore()

	s := p.Store
	heads := s.FindHeads()
	if len(heads) == 0 {
		log.Fatal("No migrations found, nothing to merge.")
	}
	if len(heads) == 1 {
		log.Fatal("There is only one migration head, there's nothing to merge.")
	}

	m := &migration.Migration{
		Name:         p.genName(),
		Dependencies: heads,
	}
	p.writeMigration(m)
}

func (p *Plugin) cmdGen(cmd *cobra.Command, args []string) {
	p.ensureStore()

	s1 := schema.NewDatabase()
	head := p.applyAll(s1)
	s2 := gen.Config.Schema.SQLSchema()
	ops := diff.Diff(s1, s2)
	if len(ops) == 0 {
		log.Fatal("No model changes found, doing nothing.")
	}

	var deps []string
	if head != "" {
		deps = []string{head}
	}

	m := &migration.Migration{
		Name:         p.genName(),
		Dependencies: deps,
		Operations:   ops,
	}
	p.writeMigration(m)
}

func (p *Plugin) genName() string {
	n := 0
	for m := range p.Store.Migrations {
		i := strings.IndexFunc(m, func(r rune) bool {
			return r < '0' || r > '9'
		})
		if i == -1 {
			i = len(m)
		}
		if i == 0 {
			continue
		}

		j, err := strconv.Atoi(m[:i])
		if err != nil {
			panic(err) // This should never happen
		}

		if n < j {
			n = j
		}
	}
	n++

	var b [3]byte
	_, _ = rand.Read(b[:])

	return fmt.Sprintf("%05d_%s", n, hex.EncodeToString(b[:]))
}

func (p *Plugin) ensureStore() {
	if err := os.MkdirAll(p.PackagePath, os.ModePerm); err != nil {
		log.Fatalf("Error creating output directory %s: %v", p.PackagePath, err)
	}

	if _, err := os.Stat(filepath.Join(p.PackagePath, "store.go")); os.IsNotExist(err) {
		var buf bytes.Buffer
		gen.WritePackageName(&buf, p.PackageName)
		buf.WriteString("import \"github.com/sqlbunny/sqlbunny/runtime/migration\"\n")
		buf.WriteString("\n")
		buf.WriteString("// Store contains the migrations for this project\n")
		buf.WriteString("var Store migration.Store\n")

		gen.WriteFile(p.PackagePath, "store.go", buf.Bytes())

		if p.Store == nil {
			log.Println("Initial migrations package created.")
			log.Println("To generate migrations, you need to add a reference to the")
			log.Println("migration store in the plugin config, like this:")
			log.Println()
			log.Println("    &migration.Plugin{")
			log.Println("        Store: &migrations.Store,")
			log.Println("    },")
			log.Println()
			p.Store = &migration.Store{}
		}
	} else {
		if p.Store == nil {
			log.Println("No migration store in the plugin config, but it seems to exist on disk!")
			log.Println("To generate migrations, you need to add a reference to the")
			log.Println("migration store in the plugin config, like this:")
			log.Println()
			log.Println("    &migration.Plugin{")
			log.Println("        Store: &migrations.Store,")
			log.Println("    },")
			log.Println()
			log.Fatal("Exiting")
		}
	}
}

func (p *Plugin) writeMigration(m *migration.Migration) {
	migrationFile := fmt.Sprintf("migration_%s.go", m.Name)
	var buf bytes.Buffer
	gen.WritePackageName(&buf, p.PackageName)
	buf.WriteString("import (\n")
	buf.WriteString("    \"github.com/sqlbunny/sqlbunny/runtime/migration\"\n")
	buf.WriteString("    \"github.com/sqlbunny/sqlschema/operations\"\n")
	buf.WriteString(")\n")
	buf.WriteString(fmt.Sprintf("func init() {\nStore.Register("))
	buf.WriteString(litter.Options{}.Sdump(m))
	buf.WriteString(")\n}")

	gen.WriteFile(p.PackagePath, migrationFile, buf.Bytes())
}

func (p *Plugin) cmdGenSQL(cmd *cobra.Command, args []string) {
	s1 := schema.NewDatabase()
	s2 := gen.Config.Schema.SQLSchema()
	ops := diff.Diff(s1, s2)
	if len(ops) == 0 {
		log.Fatal("No models found, doing nothing.")
	}

	for _, op := range ops {
		q := op.GetSQL()
		fmt.Println(q + ";\n")
	}
}

func (p *Plugin) applyAll(db *schema.Database) string {
	s := p.Store

	if len(s.Migrations) == 0 {
		return ""
	}

	heads := s.FindHeads()
	if len(heads) != 1 {
		log.Fatal("Found multiple migration heads, you must run 'migration merge' first")
	}
	head := heads[0]

	err := s.RunMigration(head, nil, func(m *migration.Migration) error {
		return ApplyMigration(m, db)
	})
	if err != nil {
		log.Fatalf("Error applying migrations: %v", err)
	}

	return head
}
