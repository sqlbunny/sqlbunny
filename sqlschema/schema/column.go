package schema

// Column holds information about a database table column.
type Column struct {
	Type     string
	Default  string
	Nullable bool
}
