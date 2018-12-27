package schema

import "github.com/kernelpayments/sqlbunny/runtime/strmangle"

type Enum struct {
	Name    string
	Choices []string
}

func (e *Enum) GetName() string {
	return e.Name
}

func (e *Enum) TypeGo() TypeGo {
	return TypeGo{
		Name: strmangle.TitleCase(e.Name),
	}
}

func (e *Enum) TypeGoNull() TypeGo {
	return TypeGo{
		Name: "Null" + strmangle.TitleCase(e.Name),
	}
}

func (e *Enum) TypeGoNullField() string {
	return strmangle.TitleCase(e.Name)
}

func (e *Enum) TypeDB() string {
	return "integer"
}
