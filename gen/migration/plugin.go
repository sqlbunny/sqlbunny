package migration

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kernelpayments/sqlbunny/runtime/bunny"

	"github.com/kernelpayments/sqlbunny/gen"
	"github.com/kernelpayments/sqlbunny/runtime/migration"
	"github.com/kernelpayments/sqlbunny/schema"
	"github.com/spf13/cobra"
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
	gen.AddCommand(&cobra.Command{
		Use: "genmigrations",
		Run: p.cmdGenMigrations,
	})
	gen.AddCommand(&cobra.Command{
		Use: "gensql",
		Run: p.cmdGenSQL,
	})
}

func (p *Plugin) migrationsOutputPath() string {
	return filepath.Join(gen.Config.OutputPath, p.MigrationsPackageName)
}

func (p *Plugin) cmdGenMigrations(cmd *cobra.Command, args []string) {
	if err := os.MkdirAll(p.migrationsOutputPath(), os.ModePerm); err != nil {
		log.Fatalf("Error creating output directory %s: %v", p.migrationsOutputPath(), err)
	}

	if _, err := os.Stat(filepath.Join(p.migrationsOutputPath(), "store.go")); os.IsNotExist(err) {
		var buf bytes.Buffer
		gen.WritePackageName(&buf, p.MigrationsPackageName)
		buf.WriteString("import \"github.com/kernelpayments/sqlbunny/runtime/migration\"\n")
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
			log.Println("Once you've done this, run genmigrations again.")
			return
		}
	}

	if p.Store == nil {
		log.Fatal("migrate.Plugin.Store is not set.")
	}

	s1 := schema.New()
	p.applyAll(s1)
	s2 := gen.Config.Schema
	ops := diff(nil, s1, s2)

	if len(ops) == 0 {
		log.Fatal("No model changes found, doing nothing.")
	}

	migrationNumber := p.nextFree()
	migrationFile := fmt.Sprintf("migration_%05d.go", migrationNumber)

	var buf bytes.Buffer
	gen.WritePackageName(&buf, p.MigrationsPackageName)
	buf.WriteString("import \"github.com/kernelpayments/sqlbunny/runtime/migration\"\n")
	buf.WriteString(fmt.Sprintf("func init() {\nStore.Register(%d, ", migrationNumber))
	ops.Dump(&buf)
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
		log.Fatal("No model changes found, doing nothing.")
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

func (p *Plugin) applyAll(db *schema.Schema) {
	var i int64 = 1
	for {
		ops, ok := p.Store.Migrations[i]
		if !ok {
			break
		}

		for _, op := range ops {
			ApplyOperation(op, db)
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
