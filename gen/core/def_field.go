package core

import (
	"fmt"

	"github.com/sqlbunny/sqlbunny/schema"
)

type StructFieldContext struct {
	*StructContext
	Field *schema.Field
}

type StructFieldItem interface {
	StructFieldItem(ctx *StructFieldContext)
}

type ModelFieldContext struct {
	*ModelContext
	Field *schema.Field
}

type ModelFieldItem interface {
	ModelFieldItem(ctx *ModelFieldContext)
}

type defModelField struct {
	name     string
	typeName string
	items    []ModelFieldItem
}

func (d defModelField) StructItem(ctx *StructContext) {
	f := &schema.Field{
		Name: d.name,
		Type: ctx.GetType(d.typeName, fmt.Sprintf("Struct %s, field %s", ctx.Struct.Name, d.name)),
		Tags: schema.Tags{},
	}

	ctx.Struct.Fields = append(ctx.Struct.Fields, f)

	for _, i := range d.items {
		if i, ok := i.(StructFieldItem); ok {
			i.StructFieldItem(&StructFieldContext{
				StructContext: ctx,
				Field:         f,
			})
		}
	}
}

func (d defModelField) ModelItem(ctx *ModelContext) {
	m := ctx.Model

	t := ctx.GetType(d.typeName, fmt.Sprintf("Model '%s' field '%s'", ctx.Model.Name, ctx.Prefix+d.name))
	if t == nil {
		return
	}

	f := &schema.Field{
		Name:     d.name,
		Type:     t,
		Nullable: false,
		Tags:     schema.Tags{},
	}
	for _, i := range d.items {
		i.ModelFieldItem(&ModelFieldContext{
			ModelContext: ctx,
			Field:        f,
		})
	}

	if ctx.Prefix == "" {
		m.Fields = append(m.Fields, f)
	}

	switch t := t.(type) {
	case *schema.Struct:
		defStruct := t.GetExtension(defStructExt{}).(*structType)

		ctx2 := &ModelContext{
			Context:       ctx.Context,
			Model:         ctx.Model,
			Prefix:        ctx.Prefix + d.name + ".",
			ForceNullable: ctx.ForceNullable || f.Nullable,
		}
		for _, i := range defStruct.items {
			if i, ok := i.(ModelItem); ok {
				i.ModelItem(ctx2)
			}
		}

		if f.Nullable {
			var def string
			if !ctx.ForceNullable {
				def = "false"
			}
			m.Columns = append(m.Columns, &schema.Column{
				Name: undot(ctx.Prefix + d.name),
				Type: &schema.BaseTypeNullable{
					Name: "bool",
					Go: schema.GoType{
						Name: "bool",
					},
					GoNull: schema.GoType{
						Pkg:  "github.com/sqlbunny/sqlbunny/types/null",
						Name: "Bool",
					},
					Postgres: schema.SQLType{
						Type:      "boolean",
						ZeroValue: "false",
					},
				},
				SQLType:    "boolean",
				SQLDefault: def,
				Nullable:   ctx.ForceNullable,
			})
		}
	case schema.BaseType:
		nullable := f.Nullable || ctx.ForceNullable
		var def string
		if !nullable {
			def = t.SQLType().ZeroValue
		}
		m.Columns = append(m.Columns, &schema.Column{
			Name:       undot(ctx.Prefix + d.name),
			Type:       t,
			SQLType:    t.SQLType().Type,
			SQLDefault: def,
			Nullable:   nullable,
		})
	default:
		// Should never happen, because all types except Struct
		// implement schema.BaseType.
		panic("unknown type")
	}
}

func Field(name string, typeName string, items ...ModelFieldItem) defModelField {
	return defModelField{
		name:     name,
		typeName: typeName,
		items:    items,
	}
}
