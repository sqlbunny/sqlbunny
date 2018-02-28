package schema

// Field holds information about a database field.
// Types are Go types, converted by TranslateFieldType.
type Column struct {
	Name   string
	Type   Type
	DBType string

	Nullable   bool
	HasDefault bool
}

// ColumnNames of the Columns.
func ColumnNames(cols []*Column) []string {
	names := make([]string, len(cols))
	for i, c := range cols {
		names[i] = c.Name
	}

	return names
}

// FilterColumnsByDefault generates the list of fields that have default values
func FilterColumnsByDefault(defaults bool, columns []*Column) []*Column {
	var cols []*Column

	for _, c := range columns {
		if defaults == c.HasDefault {
			cols = append(cols, c)
		}
	}

	return cols
}
