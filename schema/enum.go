package schema

import "github.com/KernelPay/sqlboiler/boil/strmangle"

type Enum struct {
	Name     string
	typeName string
	Type     BaseType
	Choices  []string
}

func (e *Enum) GetName() string {
	return e.Name
}

func (e *Enum) TypeGo() string {
	return strmangle.TitleCase(e.Name)
}

func (e *Enum) BaseTypeGo() string {
	return e.Type.TypeGo()
}
