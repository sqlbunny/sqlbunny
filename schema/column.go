package schema

// Field holds information about a database field.
// Types are Go types, converted by TranslateFieldType.
type Column struct {
	Name   string
	Type   Type
	DBType string

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

func (c *Column) TypeGo() TypeGo {
	if c.Nullable {
		return c.Type.(NullableType).TypeGoNull()
	} else {
		return c.Type.TypeGo()
	}
}
