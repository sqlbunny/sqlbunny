package schema

import (
	"fmt"

	"github.com/kernelpayments/sqlbunny/bunny/strmangle"
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

	ToOneRelationships  []*ToOneRelationship
	ToManyRelationships []*ToManyRelationship
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

// FindColumn by name. Returns nil if not found.
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

// CanLastInsertID checks the following:
// 1. Is there only one primary key?
// 2. Does the primary key field have a default value?
// 3. Is the primary key field type one of uintX/intX?
// If the above is all true, this model can use LastInsertId
func (t *Model) CanLastInsertID() bool {
	if t.PrimaryKey == nil || len(t.PrimaryKey.Columns) != 1 {
		return false
	}

	col := t.GetColumn(t.PrimaryKey.Columns[0])

	switch col.DBType { // TODO FIXME XXX
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
	default:
		return false
	}

	return true
}

func (t *Model) IsUniqueColumn(name string) bool {
	for _, c := range t.Uniques {
		if len(c.Columns) == 1 && c.Columns[0] == name {
			return true
		}
	}
	return false
}

func (t *Model) IsStandardModel() bool {
	for _, c := range t.Fields {
		if c.Name == "id" && c.Type.TypeGo().Name == strmangle.TitleCase(t.Name)+"ID" {
			return true
		}
	}
	return false
}
