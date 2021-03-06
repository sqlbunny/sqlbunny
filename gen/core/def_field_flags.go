package core

import "fmt"

type defFieldNull struct{}

func (d defFieldNull) FieldItem() {}
func (d defFieldNull) ModelFieldItem(ctx *ModelFieldContext) {
	ctx.Field.Nullable = true
}

func (d defFieldNull) StructFieldItem(ctx *StructFieldContext) {
	ctx.Field.Nullable = true
}

var _ FieldItem = defFieldNull{}
var _ StructFieldItem = defFieldNull{}
var _ ModelFieldItem = defFieldNull{}

var Null defFieldNull

type defFieldTag struct {
	key   string
	value string
}

func (d defFieldTag) FieldItem() {}

func (d defFieldTag) ModelFieldItem(ctx *ModelFieldContext) {
	if _, ok := ctx.Field.Tags[d.key]; ok {
		ctx.AddError("%s has duplicate tag '%s'", fmt.Sprintf("model %s field %s", ctx.Model.Name, ctx.Field.Name), d.key)
	}
	ctx.Field.Tags[d.key] = d.value
}

var _ ModelFieldItem = defFieldTag{}

func (d defFieldTag) StructFieldItem(ctx *StructFieldContext) {
	if _, ok := ctx.Field.Tags[d.key]; ok {
		ctx.AddError("%s has duplicate tag '%s'", fmt.Sprintf("struct %s field %s", ctx.Struct.Name, ctx.Field.Name), d.key)
	}
	ctx.Field.Tags[d.key] = d.value
}

var _ FieldItem = defFieldTag{}
var _ ModelFieldItem = defFieldTag{}
var _ StructFieldItem = defFieldTag{}

func Tag(key string, value string) defFieldTag {
	return defFieldTag{key: key, value: value}
}
