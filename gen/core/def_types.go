package core

import (
	"strings"

	"github.com/sqlbunny/sqlbunny/gen"
	"github.com/sqlbunny/sqlbunny/runtime/strmangle"
	"github.com/sqlbunny/sqlbunny/schema"
)

type TypeContext struct {
	*gen.Context
	Name string
}

type TypeItem interface {
	TypeItem(ctx *TypeContext) schema.Type
}

type defType struct {
	name string
	item TypeItem
}

func (t defType) ConfigItem(ctx *gen.Context) {
	if _, ok := ctx.Schema.Types[t.name]; ok {
		ctx.AddError("Type '%s' is defined multiple times", t.name)
	}
	ctx.Schema.Types[t.name] = t.item.TypeItem(&TypeContext{
		Context: ctx,
		Name:    t.name,
	})
}

var _ gen.ConfigItem = defType{}

func Type(name string, t TypeItem) defType {
	return defType{
		name: name,
		item: t,
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

func (t BaseType) TypeItem(ctx *TypeContext) schema.Type {
	if t.GoNull == "" {
		return &schema.BaseTypeNotNullable{
			Name: ctx.Name,
			Postgres: schema.SQLType{
				Type:      t.Postgres.Type,
				ZeroValue: t.Postgres.ZeroValue,
			},
			Go: parseGoType(t.Go),
		}
	}
	return &schema.BaseTypeNullable{
		Name: ctx.Name,
		Postgres: schema.SQLType{
			Type:      t.Postgres.Type,
			ZeroValue: t.Postgres.ZeroValue,
		},
		Go:     parseGoType(t.Go),
		GoNull: parseGoType(t.GoNull),
	}
}

type Enum map[int]string

func (t Enum) TypeItem(ctx *TypeContext) schema.Type {
	return &schema.Enum{
		Name:    ctx.Name,
		Choices: t,
	}
}

type array struct {
	element string
}

func (t array) TypeItem(ctx *TypeContext) schema.Type {
	return &schema.BaseTypeNotNullable{
		Name: ctx.Name,
		Postgres: schema.SQLType{
			Type:      "bytea[]",
			ZeroValue: "'{}'",
		},
		Go: schema.GoType{
			Name: strmangle.TitleCase(ctx.Name),
		},
	}
}

func Array(element string) array {
	return array{element}
}
