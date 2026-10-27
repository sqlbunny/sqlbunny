package gen

import (
	"testing"

	"github.com/sqlbunny/sqlbunny/schema"
)

func TestMapKey(t *testing.T) {
	resetImports()

	cases := []struct {
		typ      schema.GoType
		wantType string
		wantKey  string
	}{
		{schema.GoType{Name: "int64"}, "int64", "k"},
		{schema.GoType{Name: "[]byte"}, "string", "string(k)"},
		// The first imported package gets the alias _import00.
		{schema.GoType{Pkg: nullPkg, Name: "Bytes"}, "_import00.String", "_import00.NewString(string(k.Bytes), k.Valid)"},
		{schema.GoType{Pkg: nullPkg, Name: "Int64"}, "_import00.Int64", "k"},
	}
	for _, c := range cases {
		if got := templateMapKeyType(c.typ); got != c.wantType {
			t.Errorf("mapKeyType(%v): got %q want %q", c.typ, got, c.wantType)
		}
		if got := templateMapKey(c.typ, "k"); got != c.wantKey {
			t.Errorf("mapKey(%v): got %q want %q", c.typ, got, c.wantKey)
		}
	}
}

func TestCompare(t *testing.T) {
	bytea := &schema.BaseTypeNullable{
		Go:     schema.GoType{Name: "[]byte"},
		GoNull: schema.GoType{Pkg: nullPkg, Name: "Bytes"},
	}
	int64Type := &schema.BaseTypeNullable{
		Go:     schema.GoType{Name: "int64"},
		GoNull: schema.GoType{Pkg: nullPkg, Name: "Int64"},
	}
	field := func(typ schema.Type, nullable bool) *schema.Field {
		return &schema.Field{Type: typ, Nullable: nullable}
	}

	cases := []struct {
		ca, cb *schema.Field
		want   string
	}{
		{field(int64Type, false), field(int64Type, false), "a == b"},
		{field(int64Type, true), field(int64Type, false), "a.Valid && a.Int64 == b"},
		{field(int64Type, false), field(int64Type, true), "b.Valid && b.Int64 == a"},
		{field(int64Type, true), field(int64Type, true), "a == b"},
		{field(bytea, false), field(bytea, false), "bytes.Equal(a, b)"},
		{field(bytea, true), field(bytea, false), "a.Valid && bytes.Equal(a.Bytes, b)"},
		{field(bytea, false), field(bytea, true), "b.Valid && bytes.Equal(b.Bytes, a)"},
		{field(bytea, true), field(bytea, true), "a.Valid == b.Valid && bytes.Equal(a.Bytes, b.Bytes)"},
	}
	for i, c := range cases {
		if got := templateCompare("a", "b", c.ca, c.cb); got != c.want {
			t.Errorf("case %d: got %q want %q", i, got, c.want)
		}
	}
}
