package migration

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sqlbunny/sqlbunny/runtime/bunny"

	"github.com/spf13/cobra"
	"github.com/sqlbunny/sqlbunny/gen"
	"github.com/sqlbunny/sqlbunny/runtime/migration"
	"github.com/sqlbunny/sqlbunny/schema"
)

type Plugin struct {
	Store                 *migration.Store
	MigrationsPackageName string
}

var _ gen.Plugin = &Plugin{}

func (*Plugin) IsConfigItem() {}

func (p *Plugin) BunnyPlugin() {
	if p.MigrationsPackageName == "" {
		p.MigrationsPackageName = "migrations"
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

func (p *Plugin) migrationsOutputPath() string {
	return filepath.Join(gen.Config.OutputPath, p.MigrationsPackageName)
}

func (p *Plugin) cmdCheck(cmd *cobra.Command, args []string) {
	if p.Store == nil {
		log.Fatal("migrate.Plugin.Store is not set.")
	}

	s1 := schema.New()
	p.applyAll(s1)
	s2 := gen.Config.Schema
	ops := diff(nil, s1, s2)

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

	s1 := schema.New()
	head := p.applyAll(s1)
	s2 := gen.Config.Schema
	ops := diff(nil, s1, s2)
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
	if err := os.MkdirAll(p.migrationsOutputPath(), os.ModePerm); err != nil {
		log.Fatalf("Error creating output directory %s: %v", p.migrationsOutputPath(), err)
	}

	if _, err := os.Stat(filepath.Join(p.migrationsOutputPath(), "store.go")); os.IsNotExist(err) {
		var buf bytes.Buffer
		gen.WritePackageName(&buf, p.MigrationsPackageName)
		buf.WriteString("import \"github.com/sqlbunny/sqlbunny/runtime/migration\"\n")
		buf.WriteString("\n")
		buf.WriteString("// Store contains the migrations for this project\n")
		buf.WriteString("var Store migration.Store\n")

		gen.WriteFile(p.migrationsOutputPath(), "store.go", buf.Bytes())

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
	gen.WritePackageName(&buf, p.MigrationsPackageName)
	buf.WriteString("import \"github.com/sqlbunny/sqlbunny/runtime/migration\"\n")
	buf.WriteString(fmt.Sprintf("func init() {\nStore.Register("))
	m.Dump(&buf)
	buf.WriteString(")\n}")

	gen.WriteFile(p.migrationsOutputPath(), migrationFile, buf.Bytes())
}

type fakeDB struct {
	queries []string
}

func (e *fakeDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	e.queries = append(e.queries, query)
	return nil, nil
}
func (e *fakeDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	e.queries = append(e.queries, query)
	return nil, nil
}
func (e *fakeDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	e.queries = append(e.queries, query)
	return nil
}

func (p *Plugin) cmdGenSQL(cmd *cobra.Command, args []string) {
	s2 := gen.Config.Schema

	s1 := schema.New()
	ops := diff(nil, s1, s2)
	if len(ops) == 0 {
		log.Fatal("No models found, doing nothing.")
	}

	db := &fakeDB{}
	ctx := bunny.ContextWithDB(context.Background(), db)
	for _, op := range ops {
		op.Run(ctx)
	}

	for _, q := range db.queries {
		fmt.Println(q + ";\n")
	}
}

func (p *Plugin) applyAll(db *schema.Schema) string {
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
