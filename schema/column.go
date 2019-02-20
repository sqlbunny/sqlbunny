package schema

// Field holds information about a database field.
// Types are Go types, converted by TranslateFieldType.
type Column struct {
	Name string
	Type Type

	SQLType    string
	SQLDefault string

	Nullable bool
}

// ColumnNames of the Columns.
func ColumnNames(cols []*Column) []string {
	names := make([]string, len(cols))
	for i, c := range cols {
		names[i] = c.Name
	}

	return names
}

func (c *Column) GoType() GoType {
	if c.Nullable {
		return c.Type.(NullableType).GoTypeNull()
	} else {
		return c.Type.GoType()
	}
}
