package core

import "github.com/sqlbunny/sqlbunny/schema"

type ModelRelationshipContext struct {
	*ModelContext
	Relationship *schema.Relationship
}

type ModelRelationshipItem interface {
	ModelRelationshipItem(ctx *ModelRelationshipContext)
}

type defRelationship struct {
	name         string
	relationship ModelRelationshipItem
}

func Relationship(name string, relationship ModelRelationshipItem) *defRelationship {
	return &defRelationship{
		name:         name,
		relationship: relationship,
	}
}

func (d *defRelationship) ModelItem(ctx *ModelContext) {
	rel := &schema.Relationship{
		Name: d.name,
	}

	d.relationship.ModelRelationshipItem(&ModelRelationshipContext{
		ModelContext: ctx,
		Relationship: rel,
	})

	ctx.Model.Relationships = append(ctx.Model.Relationships, rel)
}

type DirectRelationship struct {
	// ForeignModel is the model name this relationship relates to.
	ForeignModel string

	// ToMany indicates this model can be related with multiple ToModel instances, not just one.
	ToMany bool

	// len(LocalFields) = len(ForeignFields)
	LocalFields   []string
	ForeignFields []string

	ForeignWhere   string
	ForeignOrderBy string
}

func (d DirectRelationship) ModelRelationshipItem(ctx *ModelRelationshipContext) {
	ctx.Relationship.ForeignModel = d.ForeignModel
	ctx.Relationship.ForeignFields = parsePathsPrefix(ctx, nil, d.ForeignFields)
	ctx.Relationship.LocalFields = parsePathsPrefix(ctx, nil, d.LocalFields)
	ctx.Relationship.ToMany = d.ToMany
	ctx.Relationship.ForeignWhere = d.ForeignWhere
	ctx.Relationship.ForeignOrderBy = d.ForeignOrderBy
}

type JoinRelationship struct {
	// ForeignModel is the model name this relationship relates to.
	ForeignModel string

	// JoinModel is the model name that joins the local and foreign models.
	JoinModel string

	// ToMany indicates this model can be related with multiple ToModel instances, not just one.
	ToMany bool

	// len(LocalFields) = len(JoinLocalFields)
	// len(JoinForeignFields) = len(ForeignFields)
	LocalFields       []string
	ForeignFields     []string
	JoinLocalFields   []string
	JoinForeignFields []string

	JoinWhere      string
	ForeignWhere   string
	ForeignOrderBy string
}

func (d JoinRelationship) ModelRelationshipItem(ctx *ModelRelationshipContext) {
	ctx.Relationship.IsJoinModel = true
	ctx.Relationship.ForeignModel = d.ForeignModel
	ctx.Relationship.JoinModel = d.JoinModel
	ctx.Relationship.ForeignFields = parsePathsPrefix(ctx, nil, d.ForeignFields)
	ctx.Relationship.LocalFields = parsePathsPrefix(ctx, nil, d.LocalFields)
	ctx.Relationship.JoinForeignFields = parsePathsPrefix(ctx, nil, d.JoinForeignFields)
	ctx.Relationship.JoinLocalFields = parsePathsPrefix(ctx, nil, d.JoinLocalFields)
	ctx.Relationship.ToMany = d.ToMany
	ctx.Relationship.JoinWhere = d.JoinWhere
	ctx.Relationship.ForeignWhere = d.ForeignWhere
	ctx.Relationship.ForeignOrderBy = d.ForeignOrderBy
}
