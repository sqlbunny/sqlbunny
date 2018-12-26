package def

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kernelpayments/sqlbunny/schema"
)

type ConfigItem interface {
	IsConfigItem()
}

type validationContext struct {
	errors      []error
	typesByName map[string]typeEntry
}

func (v *validationContext) addError(message string, args ...interface{}) {
	v.errors = append(v.errors, fmt.Errorf(message, args...))
}

func Validate(items []ConfigItem) (*Config, error) {
	s := schema.New()

	v := &validationContext{
		typesByName: make(map[string]typeEntry),
	}

	// Register all types by name
	for _, i := range items {
		t, ok := i.(typeEntry)
		if ok {
			if _, ok := s.Types[t.name]; ok {
				v.addError("Type '%s' is defined multiple times", t.name)
			}
			v.typesByName[t.name] = t
			s.Types[t.name] = t.info.getType(t.name)
		}
	}

	// Resolve all type->type references. We have to do this after all types are
	// registered, because types don't have to be declared in topological order.
	for _, t := range v.typesByName {
		t.info.resolveTypes(v, s.Types[t.name], func(name string, context string) schema.Type {
			res, ok := s.Types[name]
			if !ok {
				if context != "" {
					context = " " + context
				}
				v.addError("Type '%s'%s references unknown type '%s'", t.name, context, name)
			}
			return res
		})
	}

	// Register all models
	for _, i := range items {
		m, ok := i.(model)
		if ok {
			if _, ok := s.Models[m.name]; ok {
				v.addError("Model '%s' is defined multiple times", m.name)
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

	return &Config{
		Schema: s,
		Items:  items,
	}, nil
}

func makeTags(v *validationContext, flags []FieldFlag, context string) schema.Tags {
	res := make(schema.Tags)
	for _, i := range flags {
		if i, ok := i.(tagFlag); ok {
			if _, ok := res[i.key]; ok {
				v.addError("%s has duplicate tag '%s'", context, i.key)
			}
			res[i.key] = i.value
		}
	}
	return res
}

func isNullable(flags []FieldFlag) bool {
	for _, i := range flags {
		if _, ok := i.(nullFlag); ok {
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

func makeModel(v *validationContext, s *schema.Schema, m *schema.Model, items []ModelItem, prefix string, forceNullable bool) {
	for _, i := range items {
		switch i := i.(type) {
		case field:
			t, ok := s.Types[i.typeName]
			if !ok {
				v.addError("Model '%s' field '%s' references unknown type '%s'", m.Name, prefix+i.name, i.typeName)
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
							Go: schema.TypeGo{
								Name: "bool",
							},
							GoNull: schema.TypeGo{
								Pkg:  "github.com/kernelpayments/sqlbunny/types/null",
								Name: "Bool",
							},
							Postgres: "boolean",
						},
						DBType:   "boolean",
						Nullable: forceNullable,
					})
				}
			case schema.BaseType:
				m.Columns = append(m.Columns, &schema.Column{
					Name:     undot(prefix + i.name),
					Type:     t,
					DBType:   t.TypeDB(),
					Nullable: nullable || forceNullable,
				})
			default:
				// Should never happen, because all types except Struct
				// implement schema.BaseType.
				panic("unknown type")
			}
		case primaryKey:
			if m.PrimaryKey != nil {
				v.addError("Model '%s' has multiple primary key definitions", m.Name)
			}
			m.PrimaryKey = &schema.PrimaryKey{
				Columns: undotAll(prefixAll(i.names, prefix)),
			}
		case index:
			m.Indexes = append(m.Indexes, &schema.Index{
				Columns: undotAll(prefixAll(i.names, prefix)),
			})
		case unique:
			m.Uniques = append(m.Uniques, &schema.Unique{
				Columns: undotAll(prefixAll(i.names, prefix)),
			})
		case foreignKey:
			m.ForeignKeys = append(m.ForeignKeys, &schema.ForeignKey{
				Model:        m.Name,
				Column:       undot(prefix + i.columnName),
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

func checkDuplicateFields(v *validationContext, s *schema.Schema, m *schema.Model) {
	seen := make(map[string]struct{})
	for _, f := range m.Fields {
		if _, ok := seen[f.Name]; ok {
			v.addError("Model '%s' field '%s' is defined multiple times.", m.Name, f.Name)
		}
		seen[f.Name] = struct{}{}
	}
}

func describeIndex(columns []string) string {
	return strings.Join(columns, ", ")
}

func checkPrimaryKey(v *validationContext, s *schema.Schema, m *schema.Model) {
	pk := m.PrimaryKey

	if pk == nil {
		v.addError("Model '%s' is missing a primary key", m.Name)
	} else {
		for _, name := range pk.Columns {
			c := m.FindColumn(name)
			if c == nil {
				v.addError("Model '%s' primary key references unknown column '%s'", m.Name, name)
			} else if c.Nullable {
				v.addError("Model '%s' primary key references nullable column '%s'", m.Name, name)
			}
		}
	}
}

func checkIndexes(v *validationContext, s *schema.Schema, m *schema.Model) {
	seen := make(map[string]struct{})
	for _, f := range m.Indexes {
		f.Name = makeName(m.Name, f.Columns, "idx")

		if _, ok := seen[f.Name]; ok {
			v.addError("Model '%s' index '%s' is defined multiple times.", m.Name, describeIndex(f.Columns))
		}
		seen[f.Name] = struct{}{}

		for _, name := range f.Columns {
			c := m.FindColumn(name)
			if c == nil {
				v.addError("Model '%s' index '%s' references unknown column '%s'", m.Name, describeIndex(f.Columns), name)
			}
		}
	}
}

func checkUniques(v *validationContext, s *schema.Schema, m *schema.Model) {
	seen := make(map[string]struct{})
	for _, f := range m.Uniques {
		f.Name = makeName(m.Name, f.Columns, "key")

		if _, ok := seen[f.Name]; ok {
			v.addError("Model '%s' unique '%s' is defined multiple times.", m.Name, describeIndex(f.Columns))
		}
		seen[f.Name] = struct{}{}

		for _, name := range f.Columns {
			c := m.FindColumn(name)
			if c == nil {
				v.addError("Model '%s' unique '%s' references unknown column '%s'", m.Name, describeIndex(f.Columns), name)
			}
		}
	}
}

func checkForeignKeys(v *validationContext, s *schema.Schema, m *schema.Model) {
	for _, f := range m.ForeignKeys {
		f.Name = makeName(m.Name, []string{f.Column}, "fkey")

		c := m.FindColumn(f.Column)
		if c == nil {
			v.addError("Model '%s' has a foreign key on non-existing field '%s'", m.Name, f.Column)
		}

		m2, ok := s.Models[f.ForeignModel]
		if !ok {
			v.addError("Model '%s' field '%s' has foreign key to non-existing model '%s'", m.Name, f.Column, f.ForeignModel)
			continue
		}
		if len(m2.PrimaryKey.Columns) != 1 {
			v.addError("Model '%s' field '%s' has foreign key to model with multi-column primary key '%s'", m.Name, f.Column, f.ForeignModel)
		}
		ff := m2.PrimaryKey.Columns[0]
		f.ForeignColumn = ff
	}
}
