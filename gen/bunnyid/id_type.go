package bunnyid

import "github.com/kernelpayments/sqlbunny/runtime/strmangle"
import "github.com/kernelpayments/sqlbunny/schema"

type IDType struct {
	Name   string
	Prefix string
}

func (t *IDType) GetName() string {
	return t.Name
}
func (t *IDType) TypeGo() schema.TypeGo {
	return schema.TypeGo{
		Name: strmangle.TitleCase(t.Name),
	}
}
func (t *IDType) TypeGoNull() schema.TypeGo {
	return schema.TypeGo{
		Name: "Null" + strmangle.TitleCase(t.Name),
	}
}
func (t *IDType) TypeGoNullField() string {
	return "ID"
}
func (t *IDType) TypeDB() string {
	return "bytea"
}
