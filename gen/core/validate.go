package core

import (
	"fmt"
	"strings"

	"github.com/sqlbunny/sqlbunny/gen"
	"github.com/sqlbunny/sqlbunny/schema"
)

func buildSchema(items []gen.ConfigItem) (*schema.Schema, error) {
	ctx := &gen.Context{
		Schema: schema.New(),
	}

	for _, i := range items {
		i.ConfigItem(ctx)
	}

	ctx.Run()

	for _, m := range ctx.Schema.Models {
		checkDuplicateFields(ctx, m)
		checkPrimaryKey(ctx, m)
		checkIndexes(ctx, m)
		checkUniques(ctx, m)
		checkForeignKeys(ctx, m)
	}

	// TODO disallow double underscore.
	// TODO check FK columns match type (Go type? or just Postgres type?)

	if err := ctx.Error(); err != nil {
		return nil, err
	}

	ctx.Schema.CalculateRelationships()

	return ctx.Schema, nil
}

func undot(s string) string {
	return strings.Replace(s, ".", "__", -1)
}

func undotAll(s []string) []string {
	res := make([]string, len(s))
	for i := range s {
		res[i] = undot(s[i])
	}
	return res
}

func prefixAll(s []string, prefix string) []string {
	res := make([]string, len(s))
	for i := range s {
		res[i] = prefix + s[i]
	}
	return res
}

func makeName(model string, columns []string, suffix string) string {
	// Triple underscore because column names can have double underscores
	// if they belong to a struct.
	return fmt.Sprintf("%s___%s___%s", model, strings.Join(columns, "___"), suffix)
}

func checkDuplicateFields(ctx *gen.Context, m *schema.Model) {
	seen := make(map[string]struct{})
	for _, f := range m.Fields {
		if _, ok := seen[f.Name]; ok {
			ctx.AddError("Model '%s' field '%s' is defined multiple times.", m.Name, f.Name)
		}
		seen[f.Name] = struct{}{}
	}
}

func describeIndex(columns []string) string {
	return strings.Join(columns, ", ")
}

func checkPrimaryKey(ctx *gen.Context, m *schema.Model) {
	pk := m.PrimaryKey

	if pk == nil {
		ctx.AddError("Model '%s' is missing a primary key", m.Name)
	} else {
		for _, name := range pk.Columns {
			c := m.FindColumn(name)
			if c == nil {
				ctx.AddError("Model '%s' primary key references unknown column '%s'", m.Name, name)
			} else if c.Nullable {
				ctx.AddError("Model '%s' primary key references nullable column '%s'", m.Name, name)
			}
		}
	}
}

func checkIndexes(ctx *gen.Context, m *schema.Model) {
	seen := make(map[string]struct{})
	for _, f := range m.Indexes {
		f.Name = makeName(m.Name, f.Columns, "idx")

		if _, ok := seen[f.Name]; ok {
			ctx.AddError("Model '%s' index '%s' is defined multiple times.", m.Name, describeIndex(f.Columns))
		}
		seen[f.Name] = struct{}{}

		for _, name := range f.Columns {
			c := m.FindColumn(name)
			if c == nil {
				ctx.AddError("Model '%s' index '%s' references unknown column '%s'", m.Name, describeIndex(f.Columns), name)
			}
		}
	}
}

func checkUniques(ctx *gen.Context, m *schema.Model) {
	seen := make(map[string]struct{})
	for _, f := range m.Uniques {
		f.Name = makeName(m.Name, f.Columns, "key")

		if _, ok := seen[f.Name]; ok {
			ctx.AddError("Model '%s' unique '%s' is defined multiple times.", m.Name, describeIndex(f.Columns))
		}
		seen[f.Name] = struct{}{}

		for _, name := range f.Columns {
			c := m.FindColumn(name)
			if c == nil {
				ctx.AddError("Model '%s' unique '%s' references unknown column '%s'", m.Name, describeIndex(f.Columns), name)
			}
		}
	}
}

func checkForeignKeys(ctx *gen.Context, m *schema.Model) {
	for _, f := range m.ForeignKeys {
		f.Name = makeName(m.Name, f.LocalColumns, "fkey")

		desc := strings.Join(f.LocalColumns, ",")

		if len(f.LocalColumns) == 0 {
			ctx.AddError("Model '%s' foreign key '%s': local column list is empty", m.Name, desc)
		}
		for _, n := range f.LocalColumns {
			c := m.FindColumn(n)
			if c == nil {
				ctx.AddError("Model '%s' foreign key '%s': field '%s' does not exist", m.Name, desc, n)
			}
		}

		m2, ok := ctx.Schema.Models[f.ForeignModel]
		if !ok {
			ctx.AddError("Model '%s' foreign key '%s': foreign model '%s' does not exist", m.Name, desc, f.ForeignModel)
			continue
		}
		if f.ForeignColumns == nil && m2.PrimaryKey != nil {
			f.ForeignColumns = m2.PrimaryKey.Columns
		}

		if len(f.ForeignColumns) == 0 {
			ctx.AddError("Model '%s' foreign key '%s': foreign column list is empty", m.Name, desc)
		}
		for _, n := range f.ForeignColumns {
			fc := m2.FindColumn(n)
			if fc == nil {
				ctx.AddError("Model '%s' foreign key '%s': foreign model field '%s' does not exist", m.Name, desc, n)
			}
		}

		if len(f.LocalColumns) != len(f.ForeignColumns) {
			ctx.AddError("Model '%s' foreign key '%s': local (%d) and foreign (%d) column count doesn't match", m.Name, desc, len(f.LocalColumns), len(f.ForeignColumns))
			continue // Do not compare types if count doesn't match
		}

		for i := range f.ForeignColumns {
			fc := m2.FindColumn(f.ForeignColumns[i])
			lc := m.FindColumn(f.LocalColumns[i])
			if fc == nil || lc == nil {
				continue // Ignore these errors, they've already been reported before.
			}
			if fc.Type != lc.Type {
				ctx.AddError("Model '%s' foreign key '%s': local field '%s' and foreign field '%s' have different types: %+v %+v", m.Name, desc, f.LocalColumns[i], f.ForeignColumns[i], lc.Type, fc.Type)
			}
		}
	}
}
