package schema

type Schema struct {
	Types  map[string]Type
	Models map[string]*Model
}

func New() *Schema {
	return &Schema{
		Types:  make(map[string]Type),
		Models: make(map[string]*Model),
	}
}

func (s *Schema) CalculateRelationships() {
	for _, o := range s.Models {
		s.setIsJoinModel(o)
	}

	// Relationships have a dependency on foreign key nullability.
	for _, o := range s.Models {
		s.setForeignKeyConstraints(o)
	}

	for _, o := range s.Models {
		s.setRelationships(o)
	}
}

func (s *Schema) setIsJoinModel(t *Model) {
	if t.PrimaryKey == nil || len(t.PrimaryKey.Columns) != 2 || len(t.ForeignKeys) < 2 || len(t.Fields) > 2 {
		return
	}

	for _, c := range t.PrimaryKey.Columns {
		found := false
		for _, f := range t.ForeignKeys {
			if c == f.Column {
				found = true
				break
			}
		}
		if !found {
			return
		}
	}

	t.IsJoinModel = true
}

func (s *Schema) setForeignKeyConstraints(t *Model) {
	for i, fkey := range t.ForeignKeys {
		localColumn := t.GetColumn(fkey.Column)
		foreignModel := s.Models[fkey.ForeignModel]
		foreignColumn := foreignModel.GetColumn(fkey.ForeignColumn)

		t.ForeignKeys[i].Nullable = localColumn.Nullable
		t.ForeignKeys[i].Unique = t.IsUniqueColumn(localColumn.Name)
		t.ForeignKeys[i].ForeignColumnNullable = foreignColumn.Nullable
		t.ForeignKeys[i].ForeignColumnUnique = foreignModel.IsUniqueColumn(foreignColumn.Name)
	}
}

func (s *Schema) setRelationships(model *Model) {
	for _, t := range s.Models {
		for _, f := range t.ForeignKeys {
			if f.ForeignModel == model.Name {
				if !t.IsJoinModel && f.Unique {
					model.ToOneRelationships = append(model.ToOneRelationships, buildToOneRelationship(model, f, t))
				} else {
					model.ToManyRelationships = append(model.ToManyRelationships, buildToManyRelationship(model, f, t))
				}
			}
		}
	}
}
