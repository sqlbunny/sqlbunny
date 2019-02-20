package schema

import "github.com/kernelpayments/sqlbunny/runtime/strmangle"

type Struct struct {
	Name   string
	Fields []*Field
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
