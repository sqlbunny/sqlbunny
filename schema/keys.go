package schema

// PrimaryKey represents a primary key in a database
type PrimaryKey struct {
	Columns []string
}

// Index represents an index in a database
type Index struct {
	Columns []string
}

// Unique represents a unique constraint in a database
type Unique struct {
	Columns []string
}

// ForeignKey represents a foreign key constraint in a database
type ForeignKey struct {
	LocalColumns   []string
	ForeignTable   string
	ForeignColumns []string
}
