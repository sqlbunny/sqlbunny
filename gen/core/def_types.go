package core

import (
	"fmt"
	"strings"

	"github.com/kernelpayments/sqlbunny/runtime/strmangle"
	"github.com/kernelpayments/sqlbunny/schema"
)

type typeEntry struct {
	name string
	info TypeDef
}

func (typeEntry) IsConfigItem() {}

type TypeDef interface {
	GetType(name string) schema.Type
	ResolveTypes(v *Validation, t schema.Type, resolve func(name string, context string) schema.Type)
}

func Type(name string, t TypeDef) typeEntry {
	return typeEntry{
		name: name,
		info: t,
	}
}

type BaseType struct {
	Go       string
	GoNull   string
	Postgres SQLType
}

type SQLType struct {
	Type      string
	ZeroValue string
}

func parseGoType(s string) schema.GoType {
	i := strings.LastIndex(s, ".")
	if i == -1 {
		return schema.GoType{
			Name: s,
		}
	}
	return schema.GoType{
		Pkg:  s[:i],
		Name: s[i+1:],
	}
}

func (t BaseType) GetType(name string) schema.Type {
	if t.GoNull == "" {
		return &schema.BaseTypeNotNullable{
			Name: name,
			Postgres: schema.SQLType{
				Type:      t.Postgres.Type,
				ZeroValue: t.Postgres.ZeroValue,
			},
			Go: parseGoType(t.Go),
		}
	}
	return &schema.BaseTypeNullable{
		Name: name,
		Postgres: schema.SQLType{
			Type:      t.Postgres.Type,
			ZeroValue: t.Postgres.ZeroValue,
		},
		Go:     parseGoType(t.Go),
		GoNull: parseGoType(t.GoNull),
	}
}

func (t BaseType) ResolveTypes(v *Validation, st schema.Type, resolve func(name string, context string) schema.Type) {
}

type enum struct {
	choices []string
}

func (t enum) GetType(name string) schema.Type {
	return &schema.Enum{
		Name:    name,
		Choices: t.choices,
	}
}

func (t enum) ResolveTypes(v *Validation, st schema.Type, resolve func(name string, context string) schema.Type) {
}

func Enum(choices ...string) enum {
	return enum{choices}
}

type array struct {
	element string
}

func (t array) GetType(name string) schema.Type {
	return &schema.BaseTypeNotNullable{
		Name: name,
		Postgres: schema.SQLType{
			Type:      "bytea[]",
			ZeroValue: "'{}'",
		},
		Go: schema.GoType{
			Name: strmangle.TitleCase(name),
		},
	}
}

func (t array) ResolveTypes(v *Validation, st schema.Type, resolve func(name string, context string) schema.Type) {
}

func Array(element string) array {
	return array{element}
}

type structType struct {
	items []ModelItem
}

func (t structType) GetType(name string) schema.Type {
	return &schema.Struct{
		Name: name,
	}
}

func (t structType) ResolveTypes(v *Validation, st schema.Type, resolve func(name string, context string) schema.Type) {
	struc := st.(*schema.Struct)
	for _, i := range t.items {
		switch i := i.(type) {
		case field:
			struc.Fields = append(struc.Fields, &schema.Field{
				Name:     i.name,
				Type:     resolve(i.typeName, "field "+i.name),
				Nullable: isNullable(i.flags),
				Tags:     makeTags(v, i.flags, fmt.Sprintf("Struct '%s' field '%s'", struc.Name, i.name)),
			})
		}
	}
}

func Struct(items ...ModelItem) structType {
	return structType{
		items: expandItems(items),
	}
}
