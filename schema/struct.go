package schema

import "github.com/kernelpayments/sqlbunny/runtime/strmangle"

type Struct struct {
	Name   string
	Fields []*Field
}

func (s *Struct) GetName() string {
	return s.Name
}
func (s *Struct) TypeGo() TypeGo {
	return TypeGo{
		Name: strmangle.TitleCase(s.Name),
	}
}
func (s *Struct) TypeGoNull() TypeGo {
	return TypeGo{
		Name: "Null" + strmangle.TitleCase(s.Name),
	}
}
func (s *Struct) TypeGoNullField() string {
	return strmangle.TitleCase(s.Name)
}
