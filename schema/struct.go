package schema

import "github.com/KernelPay/sqlboiler/boil/strmangle"

type Struct struct {
	Name   string
	Fields []*Field
}

func (s *Struct) GetName() string {
	return s.Name
}
func (s *Struct) TypeGo() string {
	return strmangle.TitleCase(s.Name)
}
func (s *Struct) TypeGoNull() string {
	panic("Nullable structs not supported")
}