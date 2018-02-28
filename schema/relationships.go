package schema

// ToOneRelationship describes a relationship between two models where the local
// model has no id, and the foregin model has an id that matches a field in the
// local model, that field is also unique which changes the dynamic into a
// one-to-one style, not a to-many.
type ToOneRelationship struct {
	Model    string
	Column   string
	Nullable bool
	Unique   bool

	ForeignModel          string
	ForeignColumn         string
	ForeignColumnNullable bool
	ForeignColumnUnique   bool
}

// ToManyRelationship describes a relationship between two models where the
// local model has no id, and the foreign model has an id that matches a field
// in the local model.
type ToManyRelationship struct {
	Model    string
	Column   string
	Nullable bool
	Unique   bool

	ForeignModel          string
	ForeignColumn         string
	ForeignColumnNullable bool
	ForeignColumnUnique   bool

	ToJoinModel bool
	JoinModel   string

	JoinLocalField         string
	JoinLocalFieldNullable bool
	JoinLocalFieldUnique   bool

	JoinForeignColumn         string
	JoinForeignColumnNullable bool
	JoinForeignColumnUnique   bool
}

// ToOneRelationships relationship lookups
// Input should be the sql name of a model like: videos
func ToOneRelationships(model string, models []*Model) []*ToOneRelationship {
	localModel := GetModel(models, model)
	return toOneRelationships(localModel, models)
}

// ToManyRelationships relationship lookups
// Input should be the sql name of a model like: videos
func ToManyRelationships(model string, models []*Model) []*ToManyRelationship {
	localModel := GetModel(models, model)
	return toManyRelationships(localModel, models)
}

func toOneRelationships(model *Model, models []*Model) []*ToOneRelationship {
	var relationships []*ToOneRelationship

	for _, t := range models {
		for _, f := range t.ForeignKeys {
			if f.ForeignModel == model.Name && !t.IsJoinModel && f.Unique {
				relationships = append(relationships, buildToOneRelationship(model, f, t, models))
			}

		}
	}

	return relationships
}

func toManyRelationships(model *Model, models []*Model) []*ToManyRelationship {
	var relationships []*ToManyRelationship

	for _, t := range models {
		for _, f := range t.ForeignKeys {
			if f.ForeignModel == model.Name && (t.IsJoinModel || !f.Unique) {
				relationships = append(relationships, buildToManyRelationship(model, f, t, models))
			}
		}
	}

	return relationships
}

func buildToOneRelationship(localModel *Model, foreignKey *ForeignKey, foreignModel *Model, models []*Model) *ToOneRelationship {
	return &ToOneRelationship{
		Model:    localModel.Name,
		Column:   foreignKey.ForeignColumn,
		Nullable: foreignKey.ForeignColumnNullable,
		Unique:   foreignKey.ForeignColumnUnique,

		ForeignModel:          foreignModel.Name,
		ForeignColumn:         foreignKey.Column,
		ForeignColumnNullable: foreignKey.Nullable,
		ForeignColumnUnique:   foreignKey.Unique,
	}
}

func buildToManyRelationship(localModel *Model, foreignKey *ForeignKey, foreignModel *Model, models []*Model) *ToManyRelationship {
	if !foreignModel.IsJoinModel {
		return &ToManyRelationship{
			Model:                 localModel.Name,
			Column:                foreignKey.ForeignColumn,
			Nullable:              foreignKey.ForeignColumnNullable,
			Unique:                foreignKey.ForeignColumnUnique,
			ForeignModel:          foreignModel.Name,
			ForeignColumn:         foreignKey.Column,
			ForeignColumnNullable: foreignKey.Nullable,
			ForeignColumnUnique:   foreignKey.Unique,
			ToJoinModel:           false,
		}
	}

	relationship := &ToManyRelationship{
		Model:    localModel.Name,
		Column:   foreignKey.ForeignColumn,
		Nullable: foreignKey.ForeignColumnNullable,
		Unique:   foreignKey.ForeignColumnUnique,

		ToJoinModel: true,
		JoinModel:   foreignModel.Name,

		JoinLocalField:         foreignKey.Column,
		JoinLocalFieldNullable: foreignKey.Nullable,
		JoinLocalFieldUnique:   foreignKey.Unique,
	}

	for _, fk := range foreignModel.ForeignKeys {
		if fk == foreignKey {
			continue
		}

		relationship.JoinForeignColumn = fk.Column
		relationship.JoinForeignColumnNullable = fk.Nullable
		relationship.JoinForeignColumnUnique = fk.Unique

		relationship.ForeignModel = fk.ForeignModel
		relationship.ForeignColumn = fk.ForeignColumn
		relationship.ForeignColumnNullable = fk.ForeignColumnNullable
		relationship.ForeignColumnUnique = fk.ForeignColumnUnique
	}

	return relationship
}
