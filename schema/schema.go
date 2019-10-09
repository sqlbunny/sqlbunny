package schema

import (
	"fmt"
	"strings"

	"github.com/sqlbunny/sqlbunny/runtime/strmangle"
)

type Schema struct {
	Types  map[string]Type
	Models map[string]*Model

	Extendable
}

func New() *Schema {
	return &Schema{
		Types:  make(map[string]Type),
		Models: make(map[string]*Model),
	}
}

func (s *Schema) CalculateRelationships() {
	// Figure out which models are join models
	for _, m := range s.Models {
		m.IsJoinModel = s.isJoinModel(m)
	}

	for _, m1 := range s.Models {
		if m1.IsJoinModel {
			s.calculateJoinModelRelationships(m1)
			continue
		}

		for _, f := range m1.ForeignKeys {
			m2 := s.Models[f.ForeignModel]
			if m2.IsJoinModel {
				continue
			}

			m1.Relationships = append(m1.Relationships, &Relationship{
				Name:           localRelationshipName(f, m1, m2),
				ToMany:         false,
				IsJoinModel:    false,
				ForeignModel:   m2.Name,
				LocalColumns:   f.LocalColumns,
				ForeignColumns: f.ForeignColumns,
			})

			toMany := !m1.IsUniqueColumns(f.LocalColumns)
			m2.Relationships = append(m2.Relationships, &Relationship{
				Name:           pluralIf(foreignRelationshipName(f, m1, m2), toMany),
				ToMany:         toMany,
				IsJoinModel:    false,
				ForeignModel:   m1.Name,
				LocalColumns:   f.ForeignColumns,
				ForeignColumns: f.LocalColumns,
			})
		}
	}

	for _, m := range s.Models {
		toRemove := make(map[string]struct{})
		for i, r := range m.Relationships {
			for i2, r2 := range m.Relationships {
				if i != i2 && r.Name == r2.Name {
					toRemove[r.Name] = struct{}{}
				}
			}
		}

		for name := range toRemove {
			fmt.Printf("Warning: removing duplicate rel %s on model %s\n", name, m.Name)
			i2 := 0
			for _, r := range m.Relationships {
				if r.Name != name {
					m.Relationships[i2] = r
					i2++
				}
			}
			m.Relationships = m.Relationships[:i2]
		}
	}
}

// isJoinModel autodetects if t is a join model. A model is a join model if all are true:
// - All the columns are part of the primary key
// - There are exactly 2 foreign keys
// - The 2 foreign keys fully cover the primary key (every column belongs to one, or the other, or both)
func (s *Schema) isJoinModel(t *Model) bool {
	if t.PrimaryKey == nil || len(t.PrimaryKey.Columns) != len(t.Fields) || len(t.ForeignKeys) != 2 {
		return false
	}

	for _, c := range t.PrimaryKey.Columns {
		found := false
		for _, f := range t.ForeignKeys {
			for _, fc := range f.LocalColumns {
				if c == fc {
					found = true
				}
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func (s *Schema) calculateJoinModelRelationships(mj *Model) {
	// If m1 is a join model, we know for sure it has 2 foreign keys, so this is OK
	f1 := mj.ForeignKeys[0]
	f2 := mj.ForeignKeys[1]

	m1 := s.Models[f1.ForeignModel]
	m2 := s.Models[f2.ForeignModel]

	addJoinModelRelationship(mj, m1, m2, f1, f2)
	addJoinModelRelationship(mj, m2, m1, f2, f1)
}

func addJoinModelRelationship(mj, m1, m2 *Model, f1, f2 *ForeignKey) {
	m1.Relationships = append(m1.Relationships, &Relationship{
		Name:               strmangle.Plural(m2.Name),
		ToMany:             true,
		IsJoinModel:        true,
		JoinModel:          mj.Name,
		ForeignModel:       m2.Name,
		LocalColumns:       f1.ForeignColumns,
		JoinLocalColumns:   f1.LocalColumns,
		ForeignColumns:     f2.ForeignColumns,
		JoinForeignColumns: f2.LocalColumns,
	})
}

func pluralIf(s string, plural bool) string {
	if !plural {
		return s
	}
	return strmangle.Plural(s)
}

func localRelationshipName(f *ForeignKey, m1, m2 *Model) string {
	if len(f.LocalColumns) == 1 {
		c := f.LocalColumns[0]
		c = trimSuffixes(c)
		return clean(c)
	}
	return clean(m2.Name)
}

func foreignRelationshipName(f *ForeignKey, m1, m2 *Model) string {
	if len(f.LocalColumns) == 1 {
		c := f.LocalColumns[0]
		c = trimSuffixes(c)
		if strings.HasPrefix(m1.Name, m2.Name+"_") {
			return clean(strings.TrimPrefix(m1.Name, m2.Name+"_"))
		}
		if c == m2.Name {
			return clean(m1.Name)
		}

		if strings.HasSuffix(c, "_"+m2.Name) {
			c = strings.TrimSuffix(c, "_"+m2.Name)
		}

		return clean(c + "_" + m1.Name)
	}
	return m1.Name
}

func clean(s string) string {
	return strings.ReplaceAll(strings.Trim(s, "_"), "__", "_")
}

var identifierSuffixes = []string{"_id", "_uuid", "_guid", "_oid"}

// trimSuffixes from the identifier
func trimSuffixes(str string) string {
	ln := len(str)
	for _, s := range identifierSuffixes {
		str = strings.TrimSuffix(str, s)
		if len(str) != ln {
			break
		}
	}

	return str
}
