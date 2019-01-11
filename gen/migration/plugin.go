package migration

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/kernelpayments/sqlbunny/runtime/bunny"

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
	gen.AddCommand(&cobra.Command{
		Use: "gensql",
		Run: p.cmdGenSQL,
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

type fakeExecutor struct {
	queries []string
}

func (e *fakeExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	e.queries = append(e.queries, query)
	return nil, nil
}
func (e *fakeExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	e.queries = append(e.queries, query)
	return nil, nil
}
func (e *fakeExecutor) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
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

	e := &fakeExecutor{}
	ctx := bunny.WithExecutor(context.Background(), e)

	for _, op := range ops {
		op.Run(ctx)
	}

	for _, q := range e.queries {
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
