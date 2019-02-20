package schema

import "github.com/kernelpayments/sqlbunny/runtime/strmangle"

type Enum struct {
	Name    string
	Choices []string
}

func (e *Enum) GetName() string {
	return e.Name
}

func (e *Enum) GoType() GoType {
	return GoType{
		Name: strmangle.TitleCase(e.Name),
	}
}

func (e *Enum) GoTypeNull() GoType {
	return GoType{
		Name: "Null" + strmangle.TitleCase(e.Name),
	}
}

func (e *Enum) GoTypeNullField() string {
	return strmangle.TitleCase(e.Name)
}

func (e *Enum) SQLType() SQLType {
	return SQLType{
		Type:      "integer",
		ZeroValue: "0",
	}
}

var _ BaseType = &Enum{}
