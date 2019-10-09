package core

import "github.com/sqlbunny/sqlbunny/schema"

type defModelPrimaryKey struct {
	names []string
}

func (d defModelPrimaryKey) ModelItem(ctx *ModelContext) {
	m := ctx.Model
	if m.PrimaryKey != nil {
		ctx.AddError("Model '%s' has multiple primary key definitions", m.Name)
	}
	m.PrimaryKey = &schema.PrimaryKey{
		Columns: undotAll(prefixAll(d.names, ctx.Prefix)),
	}
}
func (d defModelPrimaryKey) StructItem(ctx *StructContext) {}

var _ ModelItem = defModelPrimaryKey{}
var _ StructItem = defModelPrimaryKey{}

type defFieldPrimaryKey func(...string) defModelPrimaryKey

func (d defFieldPrimaryKey) ModelFieldItem(ctx *ModelFieldContext) {
	m := ctx.Model
	if m.PrimaryKey != nil {
		ctx.AddError("Model '%s' has multiple primary key definitions", m.Name)
	}
	m.PrimaryKey = &schema.PrimaryKey{
		Columns: []string{undot(ctx.Prefix + ctx.Field.Name)},
	}
}

var _ ModelFieldItem = defFieldPrimaryKey(nil)

var PrimaryKey defFieldPrimaryKey = func(names ...string) defModelPrimaryKey {
	return defModelPrimaryKey{names: names}
}

type defModelIndex struct {
	names []string
}

func (d defModelIndex) ModelItem(ctx *ModelContext) {
	m := ctx.Model
	m.Indexes = append(m.Indexes, &schema.Index{
		Columns: undotAll(prefixAll(d.names, ctx.Prefix)),
	})
}
func (d defModelIndex) StructItem(ctx *StructContext) {}

var _ ModelItem = defModelIndex{}
var _ StructItem = defModelIndex{}

type defFieldIndex func(...string) defModelIndex

func (d defFieldIndex) ModelFieldItem(ctx *ModelFieldContext) {
	m := ctx.Model
	m.Indexes = append(m.Indexes, &schema.Index{
		Columns: []string{undot(ctx.Prefix + ctx.Field.Name)},
	})
}

var _ ModelFieldItem = defFieldIndex(nil)

var Index defFieldIndex = func(names ...string) defModelIndex {
	return defModelIndex{names: names}
}

type defModelUnique struct {
	names []string
}

func (d defModelUnique) ModelItem(ctx *ModelContext) {
	m := ctx.Model
	m.Uniques = append(m.Uniques, &schema.Unique{
		Columns: undotAll(prefixAll(d.names, ctx.Prefix)),
	})
}
func (d defModelUnique) StructItem(ctx *StructContext) {}

var _ ModelItem = defModelUnique{}
var _ StructItem = defModelUnique{}

type defFieldUnique func(...string) defModelUnique

func (d defFieldUnique) ModelFieldItem(ctx *ModelFieldContext) {
	m := ctx.Model
	m.Uniques = append(m.Uniques, &schema.Unique{
		Columns: []string{undot(ctx.Prefix + ctx.Field.Name)},
	})
}

var _ ModelFieldItem = defFieldUnique(nil)

var Unique defFieldUnique = func(names ...string) defModelUnique {
	return defModelUnique{names: names}
}

type defModelForeignKey struct {
	foreignModelName   string
	columnNames        []string
	foreignColumnNames []string
}

func (d defModelForeignKey) ModelItem(ctx *ModelContext) {
	m := ctx.Model
	m.ForeignKeys = append(m.ForeignKeys, &schema.ForeignKey{
		LocalColumns: undotAll(prefixAll(d.columnNames, ctx.Prefix)),
		ForeignModel: d.foreignModelName,
	})
}

var _ ModelItem = defModelForeignKey{}

func ModelForeignKey(foreignModelName string, columnNames ...string) defModelForeignKey {
	return defModelForeignKey{
		foreignModelName: foreignModelName,
		columnNames:      columnNames,
	}
}

type defFieldForeignKey struct {
	foreignModelName string
}

func (d defFieldForeignKey) ModelFieldItem(ctx *ModelFieldContext) {
	m := ctx.Model
	m.ForeignKeys = append(m.ForeignKeys, &schema.ForeignKey{
		LocalColumns: []string{undot(ctx.Prefix + ctx.Field.Name)},
		ForeignModel: d.foreignModelName,
	})
}

var _ ModelFieldItem = defFieldForeignKey{}

func ForeignKey(foreignModelName string) defFieldForeignKey {
	return defFieldForeignKey{
		foreignModelName: foreignModelName,
	}
}
