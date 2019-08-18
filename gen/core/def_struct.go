package core

import (
	"github.com/sqlbunny/sqlbunny/gen"
	"github.com/sqlbunny/sqlbunny/schema"
)

type StructContext struct {
	*gen.Context
	Struct *schema.Struct
}

type StructItem interface {
	StructItem(ctx *StructContext)
}

type structType struct {
	items []StructItem
}

type defStructExt struct{}

func (t structType) TypeItem(ctx *TypeContext) schema.Type {
	s := &schema.Struct{
		Name: ctx.Name,
	}
	s.SetExtension(defStructExt{}, &t)

	ctx.Enqueue(100, func() {
		for _, i := range t.items {
			i.StructItem(&StructContext{
				Context: ctx.Context,
				Struct:  s,
			})
		}
	})
	return s
}

func Struct(items ...StructItem) structType {
	return structType{
		items: items,
	}
}
