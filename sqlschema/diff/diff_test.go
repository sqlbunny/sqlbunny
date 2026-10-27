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
