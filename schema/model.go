package schema

import (
	"fmt"
)

// Model metadata from the database schema.
type Model struct {
	Name    string
	Fields  []*Field
	Columns []*Column

	PrimaryKey  *PrimaryKey
	Indexes     []*Index
	Uniques     []*Unique
	ForeignKeys []*ForeignKey

	IsJoinModel bool

	Relationships []*Relationship

	Extendable
}

// GetModel by name. Panics if not found (for use in templates mostly).
func GetModel(models []*Model, name string) (tbl *Model) {
	for _, t := range models {
		if t.Name == name {
			return t
		}
	}

	panic(fmt.Sprintf("could not find model name: %s", name))
}

// GetField by name. Panics if not found (for use in templates mostly).
func (t *Model) GetField(name string) (col *Field) {
	for _, c := range t.Fields {
		if c.Name == name {
			return c
		}
	}

	panic(fmt.Sprintf("could not find field name: %s", name))
}

// FindField by name. Returns nil if not found (for use in templates mostly).
func (t *Model) FindField(name string) (col *Field) {
	for _, c := range t.Fields {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// GetColumn by name. Returns nil if not found.
func (t *Model) GetColumn(name string) *Column {
	for _, c := range t.Columns {
		if c.Name == name {
			return c
		}
	}
	panic(fmt.Sprintf("could not find column name: %s", name))
}

// FindColumn by name. Returns nil if not found.
func (t *Model) FindColumn(name string) *Column {
	for _, c := range t.Columns {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// DeleteColumn by name. does nothing if not found.
func (t *Model) DeleteColumn(name string) {
	for i, c := range t.Columns {
		if c.Name == name {
			t.Columns = append(t.Columns[:i], t.Columns[i+1:]...)
			return
		}
	}
}

// FindIndex by name. Returns nil if not found.
func (t *Model) FindIndex(name string) *Index {
	for _, c := range t.Indexes {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// DeleteIndex by name. does nothing if not found.
func (t *Model) DeleteIndex(name string) {
	for i, c := range t.Indexes {
		if c.Name == name {
			t.Indexes = append(t.Indexes[:i], t.Indexes[i+1:]...)
			return
		}
	}
}

// FindUnique by name. Returns nil if not found.
func (t *Model) FindUnique(name string) *Unique {
	for _, c := range t.Uniques {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// DeleteUnique by name. does nothing if not found.
func (t *Model) DeleteUnique(name string) {
	for i, c := range t.Uniques {
		if c.Name == name {
			t.Uniques = append(t.Uniques[:i], t.Uniques[i+1:]...)
			return
		}
	}
}

// FindForeignKey by name. Returns nil if not found.
func (t *Model) FindForeignKey(name string) *ForeignKey {
	for _, c := range t.ForeignKeys {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// DeleteForeignKey by name. does nothing if not found.
func (t *Model) DeleteForeignKey(name string) {
	for i, c := range t.ForeignKeys {
		if c.Name == name {
			t.ForeignKeys = append(t.ForeignKeys[:i], t.ForeignKeys[i+1:]...)
			return
		}
	}
}

func (t *Model) IsUniqueColumns(cols []string) bool {
	if t.PrimaryKey != nil && isSubset(t.PrimaryKey.Columns, cols) {
		return true
	}
	for _, c := range t.Uniques {
		if isSubset(c.Columns, cols) {
			return true
		}
	}
	return false
}

func isSubset(a, b []string) bool {
	for _, ca := range a {
		found := false
		for _, cb := range b {
			if ca == cb {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
