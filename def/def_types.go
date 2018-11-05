package def

import (
	"fmt"
	"strings"

	"github.com/kernelpayments/sqlbunny/bunny/strmangle"
	"github.com/kernelpayments/sqlbunny/schema"
)

type typeEntry struct {
	name string
	info typeInfo
}

type typeInfo interface {
	getType(name string) schema.Type
	resolveTypes(t schema.Type, resolve func(name string, context string) schema.Type)
}

var types []typeEntry

func Type(name string, t typeInfo) {
	types = append(types, typeEntry{
		name: name,
		info: t,
	})
}

type BaseType struct {
	Go       string
	GoNull   string
	Postgres string
}

func parseGoType(s string) schema.TypeGo {
	i := strings.LastIndex(s, ".")
	if i == -1 {
		return schema.TypeGo{
			Name: s,
		}
	}
	return schema.TypeGo{
		Pkg:  s[:i],
		Name: s[i+1:],
	}
}

func (t BaseType) getType(name string) schema.Type {
	if t.GoNull == "" {
		return &schema.BaseTypeNotNullable{
			Name:     name,
			Postgres: t.Postgres,
			Go:       parseGoType(t.Go),
		}
	}
	return &schema.BaseTypeNullable{
		Name:     name,
		Postgres: t.Postgres,
		Go:       parseGoType(t.Go),
		GoNull:   parseGoType(t.GoNull),
	}
}

func (t BaseType) resolveTypes(st schema.Type, resolve func(name string, context string) schema.Type) {
}

type enum struct {
	choices []string
}

func (t enum) getType(name string) schema.Type {
	return &schema.Enum{
		Name:    name,
		Choices: t.choices,
	}
}

func (t enum) resolveTypes(st schema.Type, resolve func(name string, context string) schema.Type) {
}

func Enum(choices ...string) enum {
	return enum{choices}
}

type array struct {
	element string
}

func (t array) getType(name string) schema.Type {
	return &schema.BaseTypeNotNullable{
		Name:     name,
		Postgres: "bytea[]",
		Go: schema.TypeGo{
			Name: strmangle.TitleCase(name),
		},
	}
}

func (t array) resolveTypes(st schema.Type, resolve func(name string, context string) schema.Type) {
}

func Array(element string) array {
	return array{element}
}

type ID struct {
	Prefix string
}

func (t ID) getType(name string) schema.Type {
	return &schema.IDType{
		Name:   name,
		Prefix: t.Prefix,
	}
}

func (t ID) resolveTypes(st schema.Type, resolve func(name string, context string) schema.Type) {
}

type structType struct {
	items []ModelItem
}

func (t structType) getType(name string) schema.Type {
	return &schema.Struct{
		Name: name,
	}
}

func (t structType) resolveTypes(st schema.Type, resolve func(name string, context string) schema.Type) {
	struc := st.(*schema.Struct)
	for _, i := range t.items {
		switch i := i.(type) {
		case field:
			struc.Fields = append(struc.Fields, &schema.Field{
				Name:     i.name,
				Type:     resolve(i.typeName, "field "+i.name),
				Nullable: isNullable(i.flags),
				Tags:     makeTags(i.flags, fmt.Sprintf("Struct '%s' field '%s'", struc.Name, i.name)),
			})
		}
	}
}

func Struct(items ...ModelItem) structType {
	return structType{
		items: expandItems(items),
	}
}
