package schema

import "github.com/KernelPay/sqlboiler/boil/strmangle"

type IDType struct {
	Name   string
	Prefix string
}

func (t *IDType) GetName() string {
	return t.Name
}
func (t *IDType) TypeGo() string {
	return strmangle.TitleCase(t.Name)
}
func (t *IDType) TypeGoNull() string {
	return "Null" + strmangle.TitleCase(t.Name)
}
func (t *IDType) TypeGoNullField() string {
	return "ID"
}
func (t *IDType) TypeDB() string {
	return "bytea"
}
