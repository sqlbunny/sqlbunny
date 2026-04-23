package schema

import "github.com/sqlbunny/sqlbunny/runtime/strmangle"

// EnumArrayType is a postgres integer[] column that stores a list of values
// from an Enum type. The code generator emits a Go slice wrapper type
// (<Name> → []<ElementEnum>) with Scan/Value implementations that translate
// between int32[] and []<ElementEnum>.
//
// Element is resolved after all types are registered; callers must construct
// EnumArrayType through the core.EnumArray helper, which defers the lookup
// via gen.Context.Enqueue.
type EnumArrayType struct {
	Name    string
	Element *Enum

	Extendable
}

func (t *EnumArrayType) GetName() string {
	return t.Name
}

func (t *EnumArrayType) GoType() GoType {
	return GoType{
		Name: strmangle.TitleCase(t.Name),
	}
}

func (t *EnumArrayType) SQLType() SQLType {
	return SQLType{
		Type:      "integer[]",
		ZeroValue: "'{}'",
	}
}

var _ BaseType = &EnumArrayType{}
