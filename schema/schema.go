package schema

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

type Schema struct {
	BaseTypes []BaseType
	Structs   []*Struct
	Enums     []*Enum
	Models    []*Model
	IDTypes   []*IDType

	TypesByName     map[string]Type
	BaseTypesByName map[string]BaseType
	ModelsByName    map[string]*Model
}

func NewSchema() *Schema {
	return &Schema{
		TypesByName:     make(map[string]Type),
		BaseTypesByName: make(map[string]BaseType),
		ModelsByName:    make(map[string]*Model),
	}
}

func (s *Schema) ResolveTypes() error {
	s.TypesByName = make(map[string]Type)
	s.BaseTypesByName = make(map[string]BaseType)
	s.ModelsByName = make(map[string]*Model)

	for _, t := range s.IDTypes {
		if _, ok := s.TypesByName[t.Name]; ok {
			return fmt.Errorf("Duplicated type %s", t.Name)
		}
		s.TypesByName[t.Name] = t

		t2 := &IDArrayType{t}
		if _, ok := s.TypesByName[t2.GetName()]; ok {
			return fmt.Errorf("Duplicated type %s", t2.GetName())
		}
		s.TypesByName[t2.GetName()] = t2
	}

	for _, o := range s.BaseTypes {
		if _, ok := s.TypesByName[o.GetName()]; ok {
			return fmt.Errorf("Duplicated type %s", o.GetName())
		}
		s.TypesByName[o.GetName()] = o
		s.BaseTypesByName[o.GetName()] = o
	}

	for _, o := range s.Enums {
		if _, ok := s.TypesByName[o.Name]; ok {
			return fmt.Errorf("Duplicated type %s", o.Name)
		}
		s.TypesByName[o.Name] = o
	}

	for _, o := range s.Structs {
		if _, ok := s.TypesByName[o.Name]; ok {
			return fmt.Errorf("Duplicated type %s", o.Name)
		}
		s.TypesByName[o.Name] = o
	}

	for _, o := range s.Models {
		if _, ok := s.TypesByName[o.Name]; ok {
			return fmt.Errorf("Duplicated type %s", o.Name)
		}
		if _, ok := s.ModelsByName[o.Name]; ok {
			return fmt.Errorf("Duplicated type %s", o.Name)
		}
		s.ModelsByName[o.Name] = o
	}

	for _, o := range s.Enums {
		t, ok := s.BaseTypesByName[o.typeName]
		if !ok {
			return fmt.Errorf("Couldn't find basic type %s", o.typeName)
		}
		o.Type = t
	}
	for _, o := range s.Structs {
		for _, f := range o.Fields {
			t, ok := s.TypesByName[f.typeName]
			if !ok {
				return fmt.Errorf("Couldn't find type %s", f.typeName)
			}
			f.Type = t
		}
	}
	for _, o := range s.Models {
		for _, f := range o.Fields {
			t, ok := s.TypesByName[f.typeName]
			if !ok {
				return fmt.Errorf("Couldn't find type %s", f.typeName)
			}
			f.Type = t
		}
	}
	return nil
}

var (
	grammar = MakeGrammar()
)

func ParseSchema(rs io.ReadSeeker) (*Schema, error) {
	s, err := grammar.Parse(rs)
	if err != nil {
		return nil, err
	}
	schema := s.(*Schema)

	if err = schema.ResolveTypes(); err != nil {
		return nil, err
	}

	for _, o := range schema.Models {
		err = schema.createIndexes(o)
		if err != nil {
			return nil, err
		}
	}

	if err = schema.fillPrimaryKeys(); err != nil {
		return nil, err
	}

	for _, o := range schema.Models {
		err = schema.createForeignKeys(o)
		if err != nil {
			return nil, err
		}
	}

	schema.fillKeyNames()
	// TODO check duplicate indexes
	// TODO check primary key columns are not nullable
	// TODO check PK, uniques and FK columns exist
	// TODO check FK columns match type (Go type? or just Postgres type?)

	for _, o := range schema.Models {
		o.Columns = makeColumns(nil, o.Fields, "")
	}

	for _, o := range schema.Models {
		setIsJoinModel(o)
	}

	// Relationships have a dependency on foreign key nullability.
	for _, o := range schema.Models {
		schema.setForeignKeyConstraints(o)
	}

	for _, o := range schema.Models {
		schema.setRelationships(o)
	}

	return schema, nil
}

func makeName(model string, columns []string, suffix string) string {
	// Triple underscore because column names can have double underscores
	// if they belong to a struct.
	return fmt.Sprintf("%s___%s___%s", model, strings.Join(columns, "___"), suffix)
}

func (s *Schema) fillKeyNames() {
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

func makeColumns(cols []*Column, fields []*Field, prefix string) []*Column {
	for _, f := range fields {
		switch t := f.Type.(type) {
		case BaseType:
			cols = append(cols, &Column{
				Name:     prefix + f.Name,
				Type:     t,
				DBType:   t.TypeDB(),
				Nullable: f.Nullable,
			})
		case *Enum:
			cols = append(cols, &Column{
				Name:     prefix + f.Name,
				Type:     t,
				DBType:   t.Type.TypeDB(),
				Nullable: f.Nullable,
			})
		case *Struct:
			cols = makeColumns(cols, t.Fields, prefix+f.Name+"__")
		}
	}
	return cols
}

// checkPrimaryKeys ensures every model has a primary key field
func (s *Schema) fillPrimaryKeys() error {
	var missingPkey []string
	for _, t := range s.Models {
		for _, f := range t.Fields {
			if f.primaryKey {
				if t.PrimaryKey != nil {
					return errors.Errorf("Model %s has multiple primary keys", t.Name)
				}
				t.PrimaryKey = &PrimaryKey{
					Columns: []string{f.Name},
				}
			}
		}
		if t.PrimaryKey == nil {
			missingPkey = append(missingPkey, t.Name)
			continue
		}
	}

	if len(missingPkey) != 0 {
		return errors.Errorf("primary key missing in models (%s)", strings.Join(missingPkey, ", "))
	}

	return nil
}

func (s *Schema) createIndexes(t *Model) error {
	for _, f := range t.Fields {
		if f.unique {
			t.Uniques = append(t.Uniques, &Unique{
				Columns: []string{f.Name},
			})
		}
		if f.index {
			t.Indexes = append(t.Indexes, &Index{
				Columns: []string{f.Name},
			})
		}
	}

	return nil
}

func (s *Schema) createForeignKeys(t *Model) error {
	for _, f := range t.Fields {
		if f.foreignKey != "" {
			t2, ok := s.ModelsByName[f.foreignKey]
			if !ok {
				return errors.Errorf("Model %s field %s has foreign key to non-existing model %s", t.Name, f.Name, f.foreignKey)
			}
			if len(t2.PrimaryKey.Columns) != 1 {
				return errors.Errorf("Model %s field %s has foreign key to model %s with multi-column PK", t.Name, f.Name, f.foreignKey)
			}
			ff := t2.PrimaryKey.Columns[0]
			t.ForeignKeys = append(t.ForeignKeys, &ForeignKey{
				Model:         t.Name,
				Column:        f.Name,
				ForeignModel:  t2.Name,
				ForeignColumn: ff,
			})
		}
	}

	return nil
}

func setIsJoinModel(t *Model) {
	if t.PrimaryKey == nil || len(t.PrimaryKey.Columns) != 2 || len(t.ForeignKeys) < 2 || len(t.Fields) > 2 {
		return
	}

	for _, c := range t.PrimaryKey.Columns {
		found := false
		for _, f := range t.ForeignKeys {
			if c == f.Column {
				found = true
				break
			}
		}
		if !found {
			return
		}
	}

	t.IsJoinModel = true
}

func (s *Schema) setForeignKeyConstraints(t *Model) {
	for i, fkey := range t.ForeignKeys {
		localColumn := t.GetColumn(fkey.Column)
		foreignModel := GetModel(s.Models, fkey.ForeignModel)
		foreignColumn := foreignModel.GetColumn(fkey.ForeignColumn)

		t.ForeignKeys[i].Nullable = localColumn.Nullable
		t.ForeignKeys[i].Unique = t.IsUniqueColumn(localColumn.Name)
		t.ForeignKeys[i].ForeignColumnNullable = foreignColumn.Nullable
		t.ForeignKeys[i].ForeignColumnUnique = foreignModel.IsUniqueColumn(foreignColumn.Name)
	}
}

func (s *Schema) setRelationships(t *Model) {
	t.ToOneRelationships = toOneRelationships(t, s.Models)
	t.ToManyRelationships = toManyRelationships(t, s.Models)
}
