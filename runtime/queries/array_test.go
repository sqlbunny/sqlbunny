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
