package schema

// Relationship describes a relationship between two models
type Relationship struct {
	// Name is the relationship name. It must be unique per model.
	// For ToMany relationships it should be plural.
	Name string

	// ForeignModel is the model name this relationship relates to.
	ForeignModel string

	// ToMany indicates this model can be related with multiple ToModel instances, not just one.
	ToMany bool

	// IsJoinModel indicates this relationship is through a join model.
	// If IsJoinModel is true, len(LocalColumns) = len(JoinLocalColumns) and len(JoinForeignColumns) = len(ForeignColumns)
	// If IsJoinModel is false, len(LocalColumns) = len(ForeignColumns)
	IsJoinModel        bool
	JoinModel          string
	LocalColumns       []string
	ForeignColumns     []string
	JoinLocalColumns   []string
	JoinForeignColumns []string
}
