package bunnyid

import (
	"github.com/sqlbunny/sqlbunny/gen/core"
	"github.com/sqlbunny/sqlbunny/schema"
)

type ID struct {
	Prefix string
}

var _ core.TypeItem = ID{}

func (t ID) TypeItem(ctx *core.TypeContext) schema.Type {
	return &IDType{
		Name:   ctx.Name,
		Prefix: t.Prefix,
	}
}
