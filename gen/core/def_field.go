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

type ModelRecursiveFieldContext struct {
	*ModelRecursiveContext
	Field *schema.Field
}

type ModelRecursiveFieldItem interface {
	ModelRecursiveFieldItem(ctx *ModelRecursiveFieldContext)
}

type FieldItem interface {
	FieldItem()
}

type defField struct {
	name     string
	typeName string
	items    []FieldItem

	field *schema.Field // Filled on the ModelItem pass, used in the ModelRecursiveItem pass
}

func (d *defField) StructItem(ctx *StructContext) {
	f := &schema.Field{
		Name: d.name,
		Type: ctx.GetType(d.typeName, fmt.Sprintf("Struct %s, field %s", ctx.Struct.Name, d.name)),
		Tags: schema.Tags{},
	}

	ctx.Struct.Fields = append(ctx.Struct.Fields, f)
	d.field = f

	for _, i := range d.items {
		if i, ok := i.(StructFieldItem); ok {
			i.StructFieldItem(&StructFieldContext{
				StructContext: ctx,
				Field:         f,
			})
		}
	}
}

func (d *defField) ModelItem(ctx *ModelContext) {
	m := ctx.Model

	t := ctx.GetType(d.typeName, fmt.Sprintf("Model '%s' field '%s'", ctx.Model.Name, d.name))
	if t == nil {
		return
	}

	f := &schema.Field{
		Name:     d.name,
		Type:     t,
		Nullable: false,
		Tags:     schema.Tags{},
	}
	m.Fields = append(m.Fields, f)
	d.field = f

	for _, i := range d.items {
		if i, ok := i.(ModelFieldItem); ok {
			i.ModelFieldItem(&ModelFieldContext{
				ModelContext: ctx,
				Field:        f,
			})
		}
	}
}

func (d *defField) ModelRecursiveItem(ctx *ModelRecursiveContext) {
	f := d.field

	for _, i := range d.items {
		if i, ok := i.(ModelRecursiveFieldItem); ok {
			i.ModelRecursiveFieldItem(&ModelRecursiveFieldContext{
				ModelRecursiveContext: ctx,
				Field:                 f,
			})
		}
	}

	t := ctx.GetType(d.typeName, fmt.Sprintf("Model '%s' field '%s'", ctx.Model.Name, ctx.Prefix+d.name))
	if t == nil {
		return
	}

	if t, ok := t.(*schema.Struct); ok {
		defStruct := t.GetExtension(defStructExt{}).(*structType)

		ctx2 := &ModelRecursiveContext{
			Context:       ctx.Context,
			Model:         ctx.Model,
			Prefix:        ctx.Prefix + d.name + ".",
			ForceNullable: ctx.ForceNullable || f.Nullable,
		}
		for _, i := range defStruct.items {
			if i, ok := i.(ModelRecursiveItem); ok {
				i.ModelRecursiveItem(ctx2)
			}
		}
	}
}

func Field(name string, typeName string, items ...FieldItem) *defField {
	return &defField{
		name:     name,
		typeName: typeName,
		items:    items,
	}
}
