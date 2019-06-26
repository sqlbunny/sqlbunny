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

	JoinLocalColumn         string
	JoinLocalColumnNullable bool
	JoinLocalColumnUnique   bool

	JoinForeignColumn         string
	JoinForeignColumnNullable bool
	JoinForeignColumnUnique   bool
}

func buildToOneRelationship(localModel *Model, foreignKey *ForeignKey, foreignModel *Model) *ToOneRelationship {
	return &ToOneRelationship{
		Model:    localModel.Name,
		Column:   foreignKey.ForeignColumns[0],
		Nullable: foreignKey.ForeignColumnNullable,
		Unique:   foreignKey.ForeignColumnUnique,

		ForeignModel:          foreignModel.Name,
		ForeignColumn:         foreignKey.Columns[0],
		ForeignColumnNullable: foreignKey.Nullable,
		ForeignColumnUnique:   foreignKey.Unique,
	}
}

func buildToManyRelationship(localModel *Model, foreignKey *ForeignKey, foreignModel *Model) *ToManyRelationship {
	if !foreignModel.IsJoinModel {
		return &ToManyRelationship{
			Model:                 localModel.Name,
			Column:                foreignKey.ForeignColumns[0],
			Nullable:              foreignKey.ForeignColumnNullable,
			Unique:                foreignKey.ForeignColumnUnique,
			ForeignModel:          foreignModel.Name,
			ForeignColumn:         foreignKey.Columns[0],
			ForeignColumnNullable: foreignKey.Nullable,
			ForeignColumnUnique:   foreignKey.Unique,
			ToJoinModel:           false,
		}
	}

	relationship := &ToManyRelationship{
		Model:    localModel.Name,
		Column:   foreignKey.ForeignColumns[0],
		Nullable: foreignKey.ForeignColumnNullable,
		Unique:   foreignKey.ForeignColumnUnique,

		ToJoinModel: true,
		JoinModel:   foreignModel.Name,

		JoinLocalColumn:         foreignKey.Columns[0],
		JoinLocalColumnNullable: foreignKey.Nullable,
		JoinLocalColumnUnique:   foreignKey.Unique,
	}

	for _, fk := range foreignModel.ForeignKeys {
		if fk == foreignKey {
			continue
		}

		relationship.JoinForeignColumn = fk.Columns[0]
		relationship.JoinForeignColumnNullable = fk.Nullable
		relationship.JoinForeignColumnUnique = fk.Unique

		relationship.ForeignModel = fk.ForeignModel
		relationship.ForeignColumn = fk.ForeignColumns[0]
		relationship.ForeignColumnNullable = fk.ForeignColumnNullable
		relationship.ForeignColumnUnique = fk.ForeignColumnUnique
	}

	return relationship
}
