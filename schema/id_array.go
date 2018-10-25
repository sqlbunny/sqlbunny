package schema

import "github.com/kernelpayments/sqlbunny/bunny/strmangle"

type IDArrayType struct {
	idType *IDType
}

func (t *IDArrayType) GetName() string {
	return t.idType.Name + "_array"
}
func (t *IDArrayType) TypeGo() TypeGo {
	return TypeGo{
		Name: strmangle.TitleCase(t.GetName()),
	}
}
func (t *IDArrayType) TypeDB() string {
	return "bytea[]"
}
