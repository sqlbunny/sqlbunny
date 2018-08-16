package schema

import "github.com/KernelPay/sqlboiler/boil/strmangle"

type IDType struct {
	Name   string
	Prefix string
}

func (t *IDType) GetName() string {
	return t.Name
}
func (t *IDType) TypeGo() TypeGo {
	return TypeGo{
		Name: strmangle.TitleCase(t.Name),
	}
}
func (t *IDType) TypeGoNull() TypeGo {
	return TypeGo{
		Name: "Null" + strmangle.TitleCase(t.Name),
	}
}
func (t *IDType) TypeGoNullField() string {
	return "ID"
}
func (t *IDType) TypeDB() string {
	return "bytea"
}
