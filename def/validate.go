package def

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kernelpayments/sqlbunny/schema"
)

var validationErrors []error

func addError(message string, args ...interface{}) {
	validationErrors = append(validationErrors, fmt.Errorf(message, args...))
}

var typesByName map[string]typeEntry

func Schema() (*schema.Schema, error) {
	s := schema.NewSchema()

	typesByName = make(map[string]typeEntry)
	// Register all types by name
	for _, t := range types {
		if _, ok := s.Types[t.name]; ok {
			addError("Type '%s' is defined multiple times", t.name)
		}
		typesByName[t.name] = t
		s.Types[t.name] = t.info.getType(t.name)
	}

	// Resolve all type->type references. We have to do this after all types are
	// registered, because types don't have to be declared in topological order.
	for _, t := range types {
		t.info.resolveTypes(s.Types[t.name], func(name string, context string) schema.Type {
			res, ok := s.Types[name]
			if !ok {
				if context != "" {
					context = " " + context
				}
				addError("Type '%s'%s references unknown type '%s'", t.name, context, name)
			}
			return res
		})
	}

	// Register all models
	for _, m := range models {
		if _, ok := s.Models[m.name]; ok {
			addError("Model '%s' is defined multiple times", m.name)
		}
		model := &schema.Model{
			Name: m.name,
		}
		s.Models[m.name] = model

		makeModel(s, model, m.items, "", false)
	}

	for _, o := range s.Models {
		if o.PrimaryKey == nil {
			addError("Model '%s' is missing a primary key", o.Name)
		}
	}

	for _, o := range s.Models {
		checkForeignKeys(s, o)
	}

	fillKeyNames(s)
	// TODO disallow double underscore.
	// TODO check duplicate indexes
	// TODO check primary key columns are not nullable
	// TODO check PK, uniques and FK columns exist
	// TODO check FK columns match type (Go type? or just Postgres type?)

	s.CalculateRelationships()

	if len(validationErrors) != 0 {
		var b strings.Builder
		fmt.Fprintf(&b, "%d errors found:\n", len(validationErrors))
		for _, e := range validationErrors {
			b.WriteString(e.Error())
			b.WriteRune('\n')
		}
		return nil, errors.New(b.String())
	}
	return s, nil
}

func makeTags(flags []FieldFlag, context string) schema.Tags {
	res := make(schema.Tags)
	for _, i := range flags {
		if i, ok := i.(tagFlag); ok {
			if _, ok := res[i.key]; ok {
				addError("%s has duplicate tag '%s'", context, i.key)
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

func makeModel(s *schema.Schema, m *schema.Model, items []ModelItem, prefix string, forceNullable bool) {
	for _, i := range items {
		switch i := i.(type) {
		case field:
			t, ok := s.Types[i.typeName]
			if !ok {
				addError("Model '%s' field '%s' references unknown type '%s'", m.Name, prefix+i.name, i.typeName)
				continue
			}

			nullable := isNullable(i.flags)

			if prefix == "" {
				f := &schema.Field{
					Name:     i.name,
					Type:     t,
					Nullable: nullable || forceNullable,
					Tags:     makeTags(i.flags, fmt.Sprintf("Model '%s' field '%s'", m.Name, prefix+i.name)),
				}
				m.Fields = append(m.Fields, f)
			}

			switch t := t.(type) {
			case *schema.Struct:
				unparsedStruct := typesByName[i.typeName].info.(structType)
				makeModel(s, m, unparsedStruct.items, prefix+i.name+".", nullable || forceNullable)

				if nullable {
					m.Columns = append(m.Columns, &schema.Column{
						Name: undot(prefix + i.name),
						Type: &schema.BaseTypeNullable{
							Name: "bool",
							Go: schema.TypeGo{
								Name: "bool",
							},
							GoNull: schema.TypeGo{
								Pkg:  "github.com/KernelPay/sqlboiler/types/null",
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
				addError("Model '%s' has multiple primary key definitions", m.Name)
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

func fillKeyNames(s *schema.Schema) {
	for _, m := range s.Models {
		for _, k := range m.Indexes {
			k.Name = makeName(m.Name, k.Columns, "idx")
		}
		for _, k := range m.Uniques {
			k.Name = makeName(m.Name, k.Columns, "key")
		}
		for _, k := range m.ForeignKeys {
			k.Name = makeName(m.Name, []string{k.Column}, "fkey")
		}
	}
}

func checkForeignKeys(s *schema.Schema, m *schema.Model) {
	for _, f := range m.ForeignKeys {
		m2, ok := s.Models[f.ForeignModel]
		if !ok {
			addError("Model '%s' field '%s' has foreign key to non-existing model '%s'", m.Name, f.Column, f.ForeignModel)
			continue
		}
		if len(m2.PrimaryKey.Columns) != 1 {
			addError("Model '%s' field '%s' has foreign key to model with multi-column primary key '%s'", m.Name, f.Column, f.ForeignModel)
		}
		ff := m2.PrimaryKey.Columns[0]
		f.ForeignColumn = ff
	}
}
