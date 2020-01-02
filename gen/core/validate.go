package core

import (
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
	// TODO check FK fields match type (Go type? or just Postgres type?)

	if err := ctx.Error(); err != nil {
		return nil, err
	}

	ctx.Schema.CalculateRelationships()

	// TODO remove this
	ctx.Schema.SQLSchema()

	return ctx.Schema, nil
}

type Context interface {
	AddError(message string, args ...interface{})
}

func parseIdentifier(ctx Context, s string) {
	if s == "" {
		ctx.AddError("Invalid identifier '%s': cannot be empty")
	}
	if strings.Contains(s, "__") {
		ctx.AddError("Invalid identifier '%s': cannot contain double underscores '__'")
	}
	if strings.HasPrefix(s, "_") {
		ctx.AddError("Invalid identifier '%s': cannot start with underscore '_'")
	}
	if strings.HasSuffix(s, "_") {
		ctx.AddError("Invalid identifier '%s': cannot end with underscore '_'")
	}
}

func parsePath(ctx Context, s string) schema.Path {
	res := schema.Path(strings.Split(s, "."))
	if len(res) == 0 {
		ctx.AddError("Invalid identifier '%s': cannot be empty")
	}

	for _, p := range res {
		parseIdentifier(ctx, p)
	}

	return res
}

func parsePathPrefix(ctx Context, prefix schema.Path, s string) schema.Path {
	var res schema.Path
	res = append(res, prefix...)
	res = append(res, strings.Split(s, ".")...)

	if len(res) == 0 {
		ctx.AddError("Invalid identifier '%s': cannot be empty")
	}

	for _, p := range res {
		parseIdentifier(ctx, p)
	}

	return res
}

func appendPath(path schema.Path, s string) schema.Path {
	var res schema.Path
	res = append(res, path...)
	res = append(res, s)
	return res
}

func parsePaths(ctx Context, s []string) []schema.Path {
	res := make([]schema.Path, len(s))
	for i := range s {
		res[i] = parsePath(ctx, s[i])
	}
	return res
}

func parsePathsPrefix(ctx Context, prefix schema.Path, s []string) []schema.Path {
	res := make([]schema.Path, len(s))
	for i := range s {
		res[i] = parsePathPrefix(ctx, prefix, s[i])
	}
	return res
}

func sqlNameAll(paths []schema.Path) []string {
	res := make([]string, len(paths))
	for i := range paths {
		res[i] = paths[i].SQLName()
	}
	return res
}

func dotNameAll(paths []schema.Path) []string {
	res := make([]string, len(paths))
	for i := range paths {
		res[i] = paths[i].DotName()
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

func checkDuplicateFields(ctx *gen.Context, m *schema.Model) {
	seen := make(map[string]struct{})
	for _, f := range m.Fields {
		if _, ok := seen[f.Name]; ok {
			ctx.AddError("Model '%s' field '%s' is defined multiple times.", m.Name, f.Name)
		}
		seen[f.Name] = struct{}{}
	}
}

func describeIndex(fields []schema.Path) string {
	return strings.Join(dotNameAll(fields), ", ")
}

func checkPrimaryKey(ctx *gen.Context, m *schema.Model) {
	if m.PrimaryKey == nil {
		ctx.AddError("Model '%s' is missing a primary key", m.Name)
		return
	}

	for _, p := range m.PrimaryKey.Fields {
		f := m.FindField(p)
		if f == nil {
			ctx.AddError("Model '%s' primary key references unknown field '%s'", m.Name, p.DotName())
		} else if f.Nullable {
			ctx.AddError("Model '%s' primary key references nullable field '%s'", m.Name, p.DotName())
		}
	}
}

func checkIndexes(ctx *gen.Context, m *schema.Model) {
	seen := make(map[string]struct{})
	for _, f := range m.Indexes {
		desc := describeIndex(f.Fields)

		if _, ok := seen[desc]; ok {
			ctx.AddError("Model '%s' index '%s' is defined multiple times.", m.Name, desc)
		}
		seen[desc] = struct{}{}

		for _, path := range f.Fields {
			c := m.FindField(path)
			if c == nil {
				ctx.AddError("Model '%s' index '%s' references unknown field '%s'", m.Name, desc, path.DotName())
			}
		}
	}
}

func checkUniques(ctx *gen.Context, m *schema.Model) {
	seen := make(map[string]struct{})
	for _, f := range m.Uniques {
		desc := describeIndex(f.Fields)

		if _, ok := seen[desc]; ok {
			ctx.AddError("Model '%s' unique '%s' is defined multiple times.", m.Name, desc)
		}
		seen[desc] = struct{}{}

		for _, path := range f.Fields {
			c := m.FindField(path)
			if c == nil {
				ctx.AddError("Model '%s' unique '%s' references unknown field '%s'", m.Name, desc, path.DotName())
			}
		}
	}
}

func checkForeignKeys(ctx *gen.Context, m *schema.Model) {
	seen := make(map[string]struct{})
	for _, f := range m.ForeignKeys {
		desc := strings.Join(dotNameAll(f.LocalFields), ", ")

		if len(f.LocalFields) == 0 {
			ctx.AddError("Model '%s' foreign key '%s': local field list is empty", m.Name, desc)
		}
		for _, p := range f.LocalFields {
			if f := m.FindField(p); f == nil {
				ctx.AddError("Model '%s' foreign key '%s': local field '%s' does not exist", m.Name, desc, p.DotName())
			}
		}

		m2, ok := ctx.Schema.Models[f.ForeignModel]
		if !ok {
			ctx.AddError("Model '%s' foreign key '%s': foreign model '%s' does not exist", m.Name, desc, f.ForeignModel)
			continue
		}
		if f.ForeignFields == nil && m2.PrimaryKey != nil {
			f.ForeignFields = m2.PrimaryKey.Fields
		}

		if len(f.ForeignFields) == 0 {
			ctx.AddError("Model '%s' foreign key '%s': foreign field list is empty", m.Name, desc)
		}
		for _, p := range f.ForeignFields {
			if f := m2.FindField(p); f == nil {
				ctx.AddError("Model '%s' foreign key '%s': foreign field '%s' does not exist", m.Name, desc, p.DotName())
			}
		}

		if len(f.LocalFields) != len(f.ForeignFields) {
			ctx.AddError("Model '%s' foreign key '%s': local (%d) and foreign (%d) field count doesn't match", m.Name, desc, len(f.LocalFields), len(f.ForeignFields))
			continue // Do not compare types if count doesn't match
		}

		for i := range f.ForeignFields {
			ff := m2.FindField(f.ForeignFields[i])
			lf := m.FindField(f.LocalFields[i])
			if ff == nil || lf == nil {
				continue // Ignore these errors, they've already been reported before.
			}
			if ff.Type != lf.Type {
				ctx.AddError("Model '%s' foreign key '%s': local field '%s' and foreign field '%s' have different types: %+v %+v", m.Name, desc, f.LocalFields[i], f.ForeignFields[i], lf.Type, ff.Type)
			}
		}

		if _, ok := seen[desc]; ok {
			ctx.AddError("Model '%s' foreign key '%s' is defined multiple times.", m.Name, describeIndex(f.LocalFields))
		}
		seen[desc] = struct{}{}
	}
}
