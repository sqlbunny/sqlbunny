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

func (t *IDType) GoType() schema.GoType {
	return schema.GoType{
		Name: strmangle.TitleCase(t.Name),
	}
}

func (t *IDType) GoTypeNull() schema.GoType {
	return schema.GoType{
		Name: "Null" + strmangle.TitleCase(t.Name),
	}
}

func (t *IDType) GoTypeNullField() string {
	return "ID"
}

func (t *IDType) SQLType() schema.SQLType {
	return schema.SQLType{
		Type:      "bytea",
		ZeroValue: "'\\x000000000000'",
	}
}
