package schema

import (
	"fmt"
	"strings"

	"github.com/sqlbunny/sqlschema/schema"
)

func appendPath(path Path, s string) Path {
	var res Path
	res = append(res, path...)
	res = append(res, s)
	return res
}

func sqlNameAll(paths []Path) []string {
	res := make([]string, len(paths))
	for i := range paths {
		res[i] = paths[i].SQLName()
	}
	return res
}

func makeName(model string, columns []Path, suffix string) string {
	// Triple underscore because column names can have double underscores
	// if they belong to a struct.
	return fmt.Sprintf("%s___%s___%s", model, strings.Join(sqlNameAll(columns), "___"), suffix)
}

func (s *Schema) SQLSchema() *schema.Schema {
	q := schema.New()

	for _, m := range s.Models {
		t := schema.NewTable()
		q.Tables[m.Name] = t
		m.Table = t

		for _, f := range m.Fields {
			doCalcFields(m, t, f, false, nil)
		}

		if m.PrimaryKey != nil {
			t.PrimaryKey = &schema.PrimaryKey{
				Columns: sqlNameAll(m.PrimaryKey.Fields),
			}
		}

		for _, f := range m.Indexes {
			t.Indexes[makeName(m.Name, f.Fields, "idx")] = &schema.Index{
				Columns: sqlNameAll(f.Fields),
			}
		}

		for _, f := range m.Uniques {
			t.Uniques[makeName(m.Name, f.Fields, "key")] = &schema.Unique{
				Columns: sqlNameAll(f.Fields),
			}
		}

		for _, f := range m.ForeignKeys {
			t.ForeignKeys[makeName(m.Name, f.LocalFields, "fkey")] = &schema.ForeignKey{
				ForeignTable:   f.ForeignModel,
				LocalColumns:   sqlNameAll(f.LocalFields),
				ForeignColumns: sqlNameAll(f.ForeignFields),
			}
		}
	}

	return q
}

func doCalcFields(m *Model, t *schema.Table, f *Field, forceNullable bool, prefix Path) {
	switch ty := f.Type.(type) {
	case *Struct:
		forceNullable2 := forceNullable || f.Nullable
		prefix2 := appendPath(prefix, f.Name)

		for _, f2 := range ty.Fields {
			doCalcFields(m, t, f2, forceNullable2, prefix2)
		}

		if f.Nullable {
			var def string
			if !forceNullable {
				def = "false"
			}
			colName := prefix2.SQLName()
			t.Columns[colName] = &schema.Column{
				Type:     "boolean",
				Default:  def,
				Nullable: forceNullable,
			}
		}
	case BaseType:
		nullable := f.Nullable || forceNullable
		var def string
		if !nullable {
			def = ty.SQLType().ZeroValue
		}

		colName := appendPath(prefix, f.Name).SQLName()
		t.Columns[colName] = &schema.Column{
			Type:     ty.SQLType().Type,
			Default:  def,
			Nullable: nullable,
		}
	default:
		// Should never happen, because all types except Struct
		// implement schema.BaseType.
		panic("unknown type")
	}
}
