package core

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kernelpayments/sqlbunny/gen"
	"github.com/kernelpayments/sqlbunny/schema"
)

type Validation struct {
	errors      []error
	typesByName map[string]typeEntry
}

func (v *Validation) AddError(message string, args ...interface{}) {
	v.errors = append(v.errors, fmt.Errorf(message, args...))
}

func buildSchema(items []gen.ConfigItem) (*schema.Schema, error) {
	s := schema.New()

	v := &Validation{
		typesByName: make(map[string]typeEntry),
	}

	// Register all types by name
	for _, i := range items {
		t, ok := i.(typeEntry)
		if ok {
			if _, ok := s.Types[t.name]; ok {
				v.AddError("Type '%s' is defined multiple times", t.name)
			}
			v.typesByName[t.name] = t
			s.Types[t.name] = t.info.GetType(t.name)
		}
	}

	// Resolve all type->type references. We have to do this after all types are
	// registered, because types don't have to be declared in topological order.
	for _, t := range v.typesByName {
		t.info.ResolveTypes(v, s.Types[t.name], func(name string, context string) schema.Type {
			res, ok := s.Types[name]
			if !ok {
				if context != "" {
					context = " " + context
				}
				v.AddError("Type '%s'%s references unknown type '%s'", t.name, context, name)
			}
			return res
		})
	}

	// Register all models
	for _, i := range items {
		m, ok := i.(model)
		if ok {
			if _, ok := s.Models[m.name]; ok {
				v.AddError("Model '%s' is defined multiple times", m.name)
			}
			model := &schema.Model{
				Name: m.name,
			}
			s.Models[m.name] = model

			makeModel(v, s, model, m.items, "", false)
		}
	}

	for _, m := range s.Models {
		checkDuplicateFields(v, s, m)
		checkPrimaryKey(v, s, m)
		checkIndexes(v, s, m)
		checkUniques(v, s, m)
		checkForeignKeys(v, s, m)
	}

	// TODO disallow double underscore.
	// TODO check FK columns match type (Go type? or just Postgres type?)

	if len(v.errors) != 0 {
		var b strings.Builder
		fmt.Fprintf(&b, "%d errors found:\n", len(v.errors))
		for _, e := range v.errors {
			b.WriteString(e.Error())
			b.WriteRune('\n')
		}
		return nil, errors.New(b.String())
	}

	s.CalculateRelationships()

	return s, nil
}

func makeTags(v *Validation, flags []FieldItem, context string) schema.Tags {
	res := make(schema.Tags)
	for _, i := range flags {
		if i, ok := i.(fieldTag); ok {
			if _, ok := res[i.key]; ok {
				v.AddError("%s has duplicate tag '%s'", context, i.key)
			}
			res[i.key] = i.value
		}
	}
	return res
}

func isNullable(flags []FieldItem) bool {
	for _, i := range flags {
		if _, ok := i.(fieldNull); ok {
			return true
		}
	}
	return false
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

func makeModel(v *Validation, s *schema.Schema, m *schema.Model, items []ModelItem, prefix string, forceNullable bool) {
	for _, i := range items {
		switch i := i.(type) {
		case field:
			t, ok := s.Types[i.typeName]
			if !ok {
				v.AddError("Model '%s' field '%s' references unknown type '%s'", m.Name, prefix+i.name, i.typeName)
				continue
			}

			nullable := isNullable(i.flags)

			if prefix == "" {
				f := &schema.Field{
					Name:     i.name,
					Type:     t,
					Nullable: nullable || forceNullable,
					Tags:     makeTags(v, i.flags, fmt.Sprintf("Model '%s' field '%s'", m.Name, prefix+i.name)),
				}
				m.Fields = append(m.Fields, f)
			}

			switch t := t.(type) {
			case *schema.Struct:
				unparsedStruct := v.typesByName[i.typeName].info.(structType)
				makeModel(v, s, m, unparsedStruct.items, prefix+i.name+".", nullable || forceNullable)

				if nullable {
					m.Columns = append(m.Columns, &schema.Column{
						Name: undot(prefix + i.name),
						Type: &schema.BaseTypeNullable{
							Name: "bool",
							Go: schema.GoType{
								Name: "bool",
							},
							GoNull: schema.GoType{
								Pkg:  "github.com/kernelpayments/sqlbunny/types/null",
								Name: "Bool",
							},
							Postgres: schema.SQLType{
								Type:      "boolean",
								ZeroValue: "false",
							},
						},
						SQLType:  "boolean",
						Nullable: forceNullable,
					})
				}
			case schema.BaseType:
				null := nullable || forceNullable
				def := t.SQLType().ZeroValue
				if null {
					def = ""
				}
				m.Columns = append(m.Columns, &schema.Column{
					Name:       undot(prefix + i.name),
					Type:       t,
					SQLType:    t.SQLType().Type,
					SQLDefault: def,
					Nullable:   null,
				})
			default:
				// Should never happen, because all types except Struct
				// implement schema.BaseType.
				panic("unknown type")
			}
		case modelPrimaryKey:
			if m.PrimaryKey != nil {
				v.AddError("Model '%s' has multiple primary key definitions", m.Name)
			}
			m.PrimaryKey = &schema.PrimaryKey{
				Columns: undotAll(prefixAll(i.names, prefix)),
			}
		case modelIndex:
			m.Indexes = append(m.Indexes, &schema.Index{
				Columns: undotAll(prefixAll(i.names, prefix)),
			})
		case modelUnique:
			m.Uniques = append(m.Uniques, &schema.Unique{
				Columns: undotAll(prefixAll(i.names, prefix)),
			})
		case modelForeignKey:
			var cols []string
			for _, c := range i.columnNames {
				cols = append(cols, undot(prefix+c))
			}
			m.ForeignKeys = append(m.ForeignKeys, &schema.ForeignKey{
				Model:        m.Name,
				Columns:      cols,
				ForeignModel: i.foreignModelName,
			})
		}
	}
}

func makeName(model string, columns []string, suffix string) string {
	// Triple underscore because column names can have double underscores
	// if they belong to a struct.
	return fmt.Sprintf("%s___%s___%s", model, strings.Join(columns, "___"), suffix)
}

func checkDuplicateFields(v *Validation, s *schema.Schema, m *schema.Model) {
	seen := make(map[string]struct{})
	for _, f := range m.Fields {
		if _, ok := seen[f.Name]; ok {
			v.AddError("Model '%s' field '%s' is defined multiple times.", m.Name, f.Name)
		}
		seen[f.Name] = struct{}{}
	}
}

func describeIndex(columns []string) string {
	return strings.Join(columns, ", ")
}

func checkPrimaryKey(v *Validation, s *schema.Schema, m *schema.Model) {
	pk := m.PrimaryKey

	if pk == nil {
		v.AddError("Model '%s' is missing a primary key", m.Name)
	} else {
		for _, name := range pk.Columns {
			c := m.FindColumn(name)
			if c == nil {
				v.AddError("Model '%s' primary key references unknown column '%s'", m.Name, name)
			} else if c.Nullable {
				v.AddError("Model '%s' primary key references nullable column '%s'", m.Name, name)
			}
		}
	}
}

func checkIndexes(v *Validation, s *schema.Schema, m *schema.Model) {
	seen := make(map[string]struct{})
	for _, f := range m.Indexes {
		f.Name = makeName(m.Name, f.Columns, "idx")

		if _, ok := seen[f.Name]; ok {
			v.AddError("Model '%s' index '%s' is defined multiple times.", m.Name, describeIndex(f.Columns))
		}
		seen[f.Name] = struct{}{}

		for _, name := range f.Columns {
			c := m.FindColumn(name)
			if c == nil {
				v.AddError("Model '%s' index '%s' references unknown column '%s'", m.Name, describeIndex(f.Columns), name)
			}
		}
	}
}

func checkUniques(v *Validation, s *schema.Schema, m *schema.Model) {
	seen := make(map[string]struct{})
	for _, f := range m.Uniques {
		f.Name = makeName(m.Name, f.Columns, "key")

		if _, ok := seen[f.Name]; ok {
			v.AddError("Model '%s' unique '%s' is defined multiple times.", m.Name, describeIndex(f.Columns))
		}
		seen[f.Name] = struct{}{}

		for _, name := range f.Columns {
			c := m.FindColumn(name)
			if c == nil {
				v.AddError("Model '%s' unique '%s' references unknown column '%s'", m.Name, describeIndex(f.Columns), name)
			}
		}
	}
}

func checkForeignKeys(v *Validation, s *schema.Schema, m *schema.Model) {
	for _, f := range m.ForeignKeys {
		f.Name = makeName(m.Name, f.Columns, "fkey")

		if len(f.Columns) == 0 {
			v.AddError("Model '%s' foreign key '%s': local column list is empty", m.Name, strings.Join(f.Columns, ","))
		}
		for _, n := range f.Columns {
			c := m.FindColumn(n)
			if c == nil {
				v.AddError("Model '%s' foreign key '%s': field '%s' does not exist", m.Name, strings.Join(f.Columns, ","), n)
			}
		}

		m2, ok := s.Models[f.ForeignModel]
		if !ok {
			v.AddError("Model '%s' foreign key '%s': foreign model '%s' does not exist", m.Name, strings.Join(f.Columns, ","), f.ForeignModel)
			continue
		}
		if f.ForeignColumns == nil {
			f.ForeignColumns = m2.PrimaryKey.Columns
		}

		if len(f.ForeignColumns) == 0 {
			v.AddError("Model '%s' foreign key '%s': foreign column list is empty", m.Name, strings.Join(f.Columns, ","))
		}
		for _, n := range f.ForeignColumns {
			fc := m2.FindColumn(n)
			if fc == nil {
				v.AddError("Model '%s' foreign key '%s': foreign model field '%s' does not exist", m.Name, strings.Join(f.Columns, ","), n)
			}
		}

		if len(f.Columns) != len(f.ForeignColumns) {
			v.AddError("Model '%s' foreign key '%s': local (%d) and foreign (%d) column count doesn't match", m.Name, strings.Join(f.Columns, ","), len(f.Columns), len(f.ForeignColumns))
			continue // Do not compare types if count doesn't match
		}

		for i := range f.ForeignColumns {
			fc := m2.FindColumn(f.ForeignColumns[i])
			lc := m.FindColumn(f.Columns[i])
			if fc == nil || lc == nil {
				continue // Ignore these errors, they've already been reported before.
			}
			if fc.Type != lc.Type {
				v.AddError("Model '%s' foreign key '%s': local field '%s' and foreign field '%s' have different types: %+v %+v", m.Name, strings.Join(f.Columns, ","), f.Columns[i], f.ForeignColumns[i], lc.Type, fc.Type)
			}
		}
	}
}
