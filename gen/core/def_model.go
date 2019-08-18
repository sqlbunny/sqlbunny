package core

import (
	"github.com/sqlbunny/sqlbunny/gen"
	"github.com/sqlbunny/sqlbunny/schema"
)

type ModelContext struct {
	*gen.Context

	Model         *schema.Model
	Prefix        string
	ForceNullable bool
}

type ModelItem interface {
	ModelItem(ctx *ModelContext)
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
	})
}

func Model(name string, items ...ModelItem) gen.ConfigItem {
	return defModel{
		name:  name,
		items: items,
	}
}
