package diff

import (
	"reflect"
	"testing"

	"github.com/sqlbunny/sqlbunny/sqlschema/operations"
	"github.com/sqlbunny/sqlbunny/sqlschema/schema"
)

func newDB(tables map[string]*schema.Table) *schema.Database {
	d := schema.NewDatabase()
	s := schema.NewSchema()
	for name, t := range tables {
		s.Tables[name] = t
	}
	d.Schemas[""] = s
	return d
}

func newTable(unlogged bool) *schema.Table {
	t := schema.NewTable()
	t.Unlogged = unlogged
	t.Columns["id"] = &schema.Column{Type: "bigint", Default: "0"}
	return t
}

// diffAndReplay diffs d1 and d2, applies the operations to d1, and checks that
// d1 then has no diff with d2.
func diffAndReplay(t *testing.T, d1, d2 *schema.Database) []operations.Operation {
	t.Helper()
	ops := Diff(d1, d2)
	for _, op := range ops {
		if err := op.Apply(d1); err != nil {
			t.Fatalf("apply %#v: %v", op, err)
		}
	}
	if ops2 := Diff(d1, d2); len(ops2) != 0 {
		t.Fatalf("diff after replay is not empty: %#v", ops2)
	}
	return ops
}

func TestUnloggedCreate(t *testing.T) {
	d1 := newDB(nil)
	d2 := newDB(map[string]*schema.Table{"a": newTable(true), "b": newTable(false)})
	ops := diffAndReplay(t, d1, d2)

	if len(ops) != 2 {
		t.Fatalf("want 2 ops, got %#v", ops)
	}
	want := []string{
		"CREATE UNLOGGED TABLE \"a\" (\n    \"id\" bigint NOT NULL DEFAULT 0\n)",
		"CREATE TABLE \"b\" (\n    \"id\" bigint NOT NULL DEFAULT 0\n)",
	}
	for i, op := range ops {
		if got := op.GetSQL(true); got != want[i] {
			t.Errorf("op %d: got %q want %q", i, got, want[i])
		}
	}
	if !d1.Schemas[""].Tables["a"].Unlogged || d1.Schemas[""].Tables["b"].Unlogged {
		t.Errorf("replay did not keep the unlogged flag")
	}
}

func TestUnloggedSwitch(t *testing.T) {
	cases := []struct {
		from, to bool
		want     []operations.Operation
	}{
		{false, false, nil},
		{true, true, nil},
		{false, true, []operations.Operation{operations.AlterTable{
			TableName: "a",
			Ops:       []operations.AlterTableSuboperation{operations.AlterTableSetUnlogged{}},
		}}},
		{true, false, []operations.Operation{operations.AlterTable{
			TableName: "a",
			Ops:       []operations.AlterTableSuboperation{operations.AlterTableSetLogged{}},
		}}},
	}
	for _, c := range cases {
		d1 := newDB(map[string]*schema.Table{"a": newTable(c.from)})
		d2 := newDB(map[string]*schema.Table{"a": newTable(c.to)})
		ops := diffAndReplay(t, d1, d2)
		if !reflect.DeepEqual(ops, c.want) {
			t.Errorf("%v -> %v: got %#v want %#v", c.from, c.to, ops, c.want)
		}
	}

	op := operations.AlterTable{
		TableName: "a",
		Ops: []operations.AlterTableSuboperation{
			operations.AlterTableSetUnlogged{},
			operations.AlterTableSetLogged{},
		},
	}
	want := "ALTER TABLE \"a\"\n    SET UNLOGGED,\n    SET LOGGED"
	if got := op.GetSQL(true); got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func sqlOf(ops []operations.Operation) []string {
	var res []string
	for _, op := range ops {
		res = append(res, op.GetSQL(false))
	}
	return res
}

func checkSQL(t *testing.T, ops []operations.Operation, want []string) {
	t.Helper()
	if got := sqlOf(ops); !reflect.DeepEqual(got, want) {
		t.Fatalf("got SQL %q\nwant %q", got, want)
	}
}

func tableWithCols(cols ...string) *schema.Table {
	t := newTable(false)
	for _, c := range cols {
		t.Columns[c] = &schema.Column{Type: "bytea", Nullable: true}
	}
	return t
}

func TestStatisticsCreateDrop(t *testing.T) {
	t1 := tableWithCols("org", "site")
	t2 := tableWithCols("org", "site")
	t2.Statistics["a_stat"] = &schema.Statistics{Columns: []string{"org", "site"}, Kinds: []string{"dependencies"}}
	t2.Statistics["b_stat"] = &schema.Statistics{Columns: []string{"org", "site"}}

	ops := diffAndReplay(t, newDB(map[string]*schema.Table{"x": t1}), newDB(map[string]*schema.Table{"x": t2}))
	checkSQL(t, ops, []string{
		`CREATE STATISTICS "a_stat" (dependencies) ON "org", "site" FROM "x"`,
		`CREATE STATISTICS "b_stat" ON "org", "site" FROM "x"`,
	})

	t3 := tableWithCols("org", "site")
	ops = diffAndReplay(t, newDB(map[string]*schema.Table{"x": t2}), newDB(map[string]*schema.Table{"x": t3}))
	checkSQL(t, ops, []string{
		`DROP STATISTICS "a_stat"`,
		`DROP STATISTICS "b_stat"`,
	})
}

func TestStatisticsChangeKinds(t *testing.T) {
	t1 := tableWithCols("org", "site")
	t1.Statistics["s"] = &schema.Statistics{Columns: []string{"org", "site"}, Kinds: []string{"dependencies"}}
	t2 := tableWithCols("org", "site")
	t2.Statistics["s"] = &schema.Statistics{Columns: []string{"org", "site"}, Kinds: []string{"mcv"}}

	ops := diffAndReplay(t, newDB(map[string]*schema.Table{"x": t1}), newDB(map[string]*schema.Table{"x": t2}))
	checkSQL(t, ops, []string{
		`DROP STATISTICS "s"`,
		`CREATE STATISTICS "s" (mcv) ON "org", "site" FROM "x"`,
	})
}

// Postgres drops a statistics object together with its columns, so the drop
// must come before the column drop.
func TestStatisticsDroppedBeforeColumn(t *testing.T) {
	t1 := tableWithCols("org", "site")
	t1.Statistics["s"] = &schema.Statistics{Columns: []string{"org", "site"}}
	t2 := tableWithCols("org")

	ops := diffAndReplay(t, newDB(map[string]*schema.Table{"x": t1}), newDB(map[string]*schema.Table{"x": t2}))
	checkSQL(t, ops, []string{
		`DROP STATISTICS "s"`,
		"ALTER TABLE \"x\"\n    DROP COLUMN \"site\"",
	})
}

func TestStatisticsOnNewTable(t *testing.T) {
	t2 := tableWithCols("org", "site")
	t2.Statistics["s"] = &schema.Statistics{Columns: []string{"org", "site"}}

	ops := diffAndReplay(t, newDB(nil), newDB(map[string]*schema.Table{"x": t2}))
	if len(ops) != 2 {
		t.Fatalf("want 2 ops, got %#v", ops)
	}
	if _, ok := ops[1].(operations.CreateStatistics); !ok {
		t.Fatalf("want CreateStatistics after CreateTable, got %#v", ops)
	}
}

func TestStorageParams(t *testing.T) {
	t1 := newTable(false)
	t1.StorageParams["fillfactor"] = "90"
	t1.StorageParams["autovacuum_enabled"] = "false"
	t2 := newTable(false)
	t2.StorageParams["fillfactor"] = "80"
	t2.StorageParams["autovacuum_analyze_scale_factor"] = "0.01"

	ops := diffAndReplay(t, newDB(map[string]*schema.Table{"x": t1}), newDB(map[string]*schema.Table{"x": t2}))
	checkSQL(t, ops, []string{
		"ALTER TABLE \"x\"\n    SET (autovacuum_analyze_scale_factor = 0.01),\n    SET (fillfactor = 80),\n    RESET (autovacuum_enabled)",
	})
}

func TestStorageParamsOnNewTable(t *testing.T) {
	t2 := newTable(false)
	t2.StorageParams["fillfactor"] = "90"
	t2.StorageParams["autovacuum_analyze_scale_factor"] = "0.01"

	ops := diffAndReplay(t, newDB(nil), newDB(map[string]*schema.Table{"x": t2, "y": newTable(false)}))
	checkSQL(t, ops, []string{
		"CREATE TABLE \"x\" (\n    \"id\" bigint NOT NULL DEFAULT 0\n) WITH (autovacuum_analyze_scale_factor = 0.01, fillfactor = 90)",
		"CREATE TABLE \"y\" (\n    \"id\" bigint NOT NULL DEFAULT 0\n)",
	})
	if ops[1].(operations.CreateTable).StorageParams != nil {
		t.Fatalf("want nil StorageParams for a table without them")
	}
}
