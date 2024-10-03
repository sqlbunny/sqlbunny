package core

import "github.com/sqlbunny/sqlbunny/schema"

type defModelPrimaryKey struct {
	names []string
}

func (d defModelPrimaryKey) ModelItem(ctx *ModelContext)   {}
func (d defModelPrimaryKey) StructItem(ctx *StructContext) {}
func (d defModelPrimaryKey) ModelRecursiveItem(ctx *ModelRecursiveContext) {
	m := ctx.Model
	if m.PrimaryKey != nil {
		ctx.AddError("Model '%s' has multiple primary key definitions", m.Name)
	}
	m.PrimaryKey = &schema.PrimaryKey{
		Fields: parsePathsPrefix(ctx, ctx.Prefix, d.names),
	}
}

var _ ModelItem = defModelPrimaryKey{}
var _ StructItem = defModelPrimaryKey{}
var _ ModelRecursiveItem = defModelPrimaryKey{}

type defFieldPrimaryKey func(...string) defModelPrimaryKey

func (d defFieldPrimaryKey) FieldItem() {}
func (d defFieldPrimaryKey) ModelRecursiveFieldItem(ctx *ModelRecursiveFieldContext) {
	m := ctx.Model
	if m.PrimaryKey != nil {
		ctx.AddError("Model '%s' has multiple primary key definitions", m.Name)
	}
	m.PrimaryKey = &schema.PrimaryKey{
		Fields: []schema.Path{parsePathPrefix(ctx, ctx.Prefix, ctx.Field.Name)},
	}
}

var _ FieldItem = defFieldPrimaryKey(nil)
var _ ModelRecursiveFieldItem = defFieldPrimaryKey(nil)

var PrimaryKey defFieldPrimaryKey = func(names ...string) defModelPrimaryKey {
	return defModelPrimaryKey{names: names}
}

type defModelIndex struct {
	names  []string
	where  string
	method string
}

func (d defModelIndex) Where(val string) defModelIndex {
	d.where = val
	return d
}

func (d defModelIndex) Method(val string) defModelIndex {
	d.method = val
	return d
}

func (d defModelIndex) ModelItem(ctx *ModelContext)   {}
func (d defModelIndex) StructItem(ctx *StructContext) {}

func (d defModelIndex) ModelRecursiveItem(ctx *ModelRecursiveContext) {
	m := ctx.Model
	m.Indexes = append(m.Indexes, &schema.Index{
		Fields: parsePathsPrefix(ctx, ctx.Prefix, d.names),
		Method: d.method,
		Where:  d.where,
	})
}

var _ ModelItem = defModelIndex{}
var _ StructItem = defModelIndex{}
var _ ModelRecursiveItem = defModelIndex{}

type defFieldIndex func(...string) defModelIndex

func (d defFieldIndex) FieldItem() {}

func (d defFieldIndex) ModelRecursiveFieldItem(ctx *ModelRecursiveFieldContext) {
	m := ctx.Model
	m.Indexes = append(m.Indexes, &schema.Index{
		Fields: []schema.Path{parsePathPrefix(ctx, ctx.Prefix, ctx.Field.Name)},
	})
}

var _ FieldItem = defFieldIndex(nil)
var _ ModelRecursiveFieldItem = defFieldIndex(nil)

var Index defFieldIndex = func(names ...string) defModelIndex {
	return defModelIndex{names: names}
}

type defModelUnique struct {
	names []string
}

func (d defModelUnique) ModelItem(ctx *ModelContext)   {}
func (d defModelUnique) StructItem(ctx *StructContext) {}
func (d defModelUnique) ModelRecursiveItem(ctx *ModelRecursiveContext) {
	m := ctx.Model
	m.Uniques = append(m.Uniques, &schema.Unique{
		Fields: parsePathsPrefix(ctx, ctx.Prefix, d.names),
	})
}

var _ ModelItem = defModelUnique{}
var _ StructItem = defModelUnique{}
var _ ModelRecursiveItem = defModelUnique{}

type defFieldUnique func(...string) defModelUnique

func (d defFieldUnique) FieldItem() {}
func (d defFieldUnique) ModelRecursiveFieldItem(ctx *ModelRecursiveFieldContext) {
	m := ctx.Model
	m.Uniques = append(m.Uniques, &schema.Unique{
		Fields: []schema.Path{parsePathPrefix(ctx, ctx.Prefix, ctx.Field.Name)},
	})
}

var _ FieldItem = defFieldUnique(nil)
var _ ModelRecursiveFieldItem = defFieldUnique(nil)

var Unique defFieldUnique = func(names ...string) defModelUnique {
	return defModelUnique{names: names}
}

type defModelForeignKey struct {
	foreignModelName   string
	columnNames        []string
	foreignColumnNames []string
}

func (d defModelForeignKey) ModelItem(ctx *ModelContext) {}
func (d defModelForeignKey) ModelRecursiveItem(ctx *ModelRecursiveContext) {
	m := ctx.Model
	m.ForeignKeys = append(m.ForeignKeys, &schema.ForeignKey{
		LocalFields:  parsePathsPrefix(ctx, ctx.Prefix, d.columnNames),
		ForeignModel: d.foreignModelName,
	})
}

var _ ModelItem = defModelForeignKey{}
var _ ModelRecursiveItem = defModelForeignKey{}

func ModelForeignKey(foreignModelName string, columnNames ...string) defModelForeignKey {
	return defModelForeignKey{
		foreignModelName:   foreignModelName,
		columnNames:        columnNames,
		foreignColumnNames: nil, // Autofill with the foreign model's primary key
	}
}

type defFieldForeignKey struct {
	foreignModelName string
}

func (d defFieldForeignKey) FieldItem() {}
func (d defFieldForeignKey) ModelRecursiveFieldItem(ctx *ModelRecursiveFieldContext) {
	m := ctx.Model
	m.ForeignKeys = append(m.ForeignKeys, &schema.ForeignKey{
		LocalFields:  []schema.Path{parsePathPrefix(ctx, ctx.Prefix, ctx.Field.Name)},
		ForeignModel: d.foreignModelName,
	})
}

var _ FieldItem = defFieldForeignKey{}
var _ ModelRecursiveFieldItem = defFieldForeignKey{}

func ForeignKey(foreignModelName string) defFieldForeignKey {
	return defFieldForeignKey{
		foreignModelName: foreignModelName,
	}
}
