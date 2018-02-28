package schema

import "reflect"

// Field holds information about a database field.
// Types are Go types, converted by TranslateFieldType.
type Field struct {
	Name string

	Type Type
	Tag  string

	Nullable bool

	// Temp vars for use while parsing
	typeName   string
	index      bool
	unique     bool
	primaryKey bool
	foreignKey string
}

func (f *Field) HasTag(tag string) bool {
	z := reflect.StructTag(f.Tag)
	_, ok := z.Lookup(tag)
	return ok
}

func (f *Field) IsStruct() bool {
	_, ok := f.Type.(*Struct)
	return ok
}

func (f *Field) TypeGo() string {
	if f.Nullable {
		return f.Type.(NullableType).TypeGoNull()
	} else {
		return f.Type.TypeGo()
	}
}

// FieldNames of the fields.
func FieldNames(fields []*Field) []string {
	names := make([]string, len(fields))
	for i, c := range fields {
		names[i] = c.Name
	}

	return names
}
