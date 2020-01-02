package schema

import (
	"github.com/sqlbunny/sqlbunny/runtime/strmangle"
)

type Struct struct {
	Name   string
	Fields []*Field

	Extendable
}

func (s *Struct) GetName() string {
	return s.Name
}
func (s *Struct) GoType() GoType {
	return GoType{
		Name: strmangle.TitleCase(s.Name),
	}
}
func (s *Struct) GoTypeNull() GoType {
	return GoType{
		Name: "Null" + strmangle.TitleCase(s.Name),
	}
}
func (s *Struct) GoTypeNullField() string {
	return strmangle.TitleCase(s.Name)
}

var _ Type = &Struct{}

func (s *Struct) fieldByName(name string) (col *Field) {
	for _, f := range s.Fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}
