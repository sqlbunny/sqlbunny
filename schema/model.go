package schema

import "github.com/sqlbunny/sqlschema/schema"

// Model metadata from the database schema.
type Model struct {
	Name   string
	Fields []*Field

	PrimaryKey  *PrimaryKey
	Indexes     []*Index
	Uniques     []*Unique
	ForeignKeys []*ForeignKey

	IsJoinModel bool

	Relationships []*Relationship

	Table *schema.Table

	Extendable
}

// FindField by path. Returns nil if not found.
func (m *Model) FindField(path Path) *Field {
	if len(path) == 0 {
		return nil
	}

	f := m.fieldByName(path[0])

	for _, name := range path[1:] {
		if f == nil {
			return nil
		}

		s, ok := f.Type.(*Struct)
		if !ok {
			return nil
		}
		f = s.fieldByName(name)
	}

	return f
}

func (m *Model) fieldByName(name string) *Field {
	for _, f := range m.Fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func (m *Model) IsFieldsUnique(fields []Path) bool {
	if m.PrimaryKey != nil && isSubset(m.PrimaryKey.Fields, fields) {
		return true
	}
	for _, c := range m.Uniques {
		if isSubset(c.Fields, fields) {
			return true
		}
	}
	return false
}

func isSubset(a, b []Path) bool {
	for _, ca := range a {
		found := false
		for _, cb := range b {
			if ca.Equals(cb) {
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
