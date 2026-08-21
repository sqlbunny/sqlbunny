package queries

import (
	"database/sql/driver"
	"testing"
)

type byteValuer [2]byte

func (v byteValuer) Value() (driver.Value, error) { return v[:], nil }

func TestArray(t *testing.T) {
	cases := []struct {
		got  driver.Valuer
		want string
	}{
		{Array([]int{1, 2}), "{1,2}"},
		{Array([]string{`a"b`, `c\d`}), `{"a\"b","c\\d"}`},
		{Array([]byteValuer{{0xab, 0xcd}}), `{"\\xabcd"}`},
		{Array([]*int{nil}), "{NULL}"},
		{Array([]bool{}), "{}"},
	}
	for _, c := range cases {
		v, err := c.got.Value()
		if err != nil {
			t.Fatal(err)
		}
		if v != c.want {
			t.Errorf("got %q want %q", v, c.want)
		}
	}
}

func TestBulkValues(t *testing.T) {
	sql, args := BulkValues([]string{"bytea", "text[]", "jsonb"}, [][]any{
		{byteValuer{1, 2}, "{1,2}", []byte(`{"a":1}`)},
		{byteValuer{3, 4}, "{}", nil},
	})
	want := "SELECT u.c0, u.c1::text[], u.c2 FROM unnest($1::bytea[], $2::text[], $3::jsonb[]) AS u(c0, c1, c2)"
	if sql != want {
		t.Errorf("got %q want %q", sql, want)
	}
	wantArgs := []string{`{"\\x0102","\\x0304"}`, `{"{1,2}","{}"}`, `{"{\"a\":1}",NULL}`}
	for i, a := range args {
		v, err := a.(driver.Valuer).Value()
		if err != nil || v != wantArgs[i] {
			t.Errorf("arg %d: got %q want %q (%v)", i, v, wantArgs[i], err)
		}
	}
}
