package core

import (
	"github.com/sqlbunny/sqlbunny/gen"
	"github.com/sqlbunny/sqlbunny/schema"
)

type ModelContext struct {
	*gen.Context

	Model *schema.Model
}

type ModelRecursiveContext struct {
	*gen.Context

	Model         *schema.Model
	Prefix        schema.Path
	ForceNullable bool
}

type ModelItem interface {
	ModelItem(ctx *ModelContext)
}

type ModelRecursiveItem interface {
	ModelRecursiveItem(ctx *ModelRecursiveContext)
}

type defModel struct {
	name  string
	items []ModelItem
}

func (d defModel) ConfigItem(ctx *gen.Context) {
	ctx.Enqueue(200, func() {
		if _, ok := ctx.Schema.Models[d.name]; ok {
			ctx.AddError("Model '%s' is defined multiple times", d.name)
		}
		model := &schema.Model{
			Name: d.name,
		}
		ctx.Schema.Models[d.name] = model

		for _, i := range d.items {
			i.ModelItem(&ModelContext{
				Context: ctx,
				Model:   model,
			})
		}

		ctx.Enqueue(300, func() {
			for _, i := range d.items {
				if i, ok := i.(ModelRecursiveItem); ok {
					i.ModelRecursiveItem(&ModelRecursiveContext{
						Context: ctx,
						Model:   model,
					})
				}
			}
		})
	})
}

func Model(name string, items ...ModelItem) gen.ConfigItem {
	return defModel{
		name:  name,
		items: items,
	}
}

type defModelUnlogged struct{}

func (d defModelUnlogged) ModelItem(ctx *ModelContext) {
	ctx.Model.Unlogged = true
}

var _ ModelItem = defModelUnlogged{}

// Unlogged makes the model's table a Postgres UNLOGGED table. Writes to it
// skip the WAL, so they are faster, but a crash empties the table and it is
// not copied to replicas. Use it only for data you can lose.
//
//	Model("rate_limit",
//	    Unlogged,
//	    Field("key", "bytea", PrimaryKey),
//	)
//
// A logged model cannot have a foreign key to an unlogged model.
var Unlogged defModelUnlogged
