package schema

// PrimaryKey represents a primary key in a database
type PrimaryKey struct {
	Fields []Path
}

// Index represents an index in a database
type Index struct {
	Fields []Path
	Method string // Index method. If empty, default is btree.
	Where  string // Index where clause, for partial indexes. If empty, no where clause is in effect.
}

// Unique represents a unique constraint in a database
type Unique struct {
	Fields []Path
}

// ForeignKey represents a foreign key constraint in a database
type ForeignKey struct {
	LocalFields   []Path
	ForeignModel  string
	ForeignFields []Path
}
