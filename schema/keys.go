package schema

// PrimaryKey represents a primary key in a database
type PrimaryKey struct {
	Columns []string
}

// Index represents an index in a database
type Index struct {
	Columns []string // Index columns. Order matters.
	Method  string   // Index method. If empty, default is btree.
	Where   string   // Index where clause, for partial indexes. If empty, no where clause is in effect.
}

// Unique represents a unique constraint in a database
type Unique struct {
	Columns []string
}

// ForeignKey represents a foreign key constraint in a database
type ForeignKey struct {
	LocalColumns   []string
	ForeignSchema  string
	ForeignTable   string
	ForeignColumns []string
}
