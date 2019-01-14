package schema

// Field holds information about a database field.
// Types are Go types, converted by TranslateFieldType.
type Field struct {
	Name     string
	Type     Type
	Nullable bool

	Tags Tags
}

func (f *Field) GenerateTags() string {
	if _, ok := f.Tags["bunny"]; !ok {
		f.Tags["bunny"] = f.Name
		if f.IsStruct() {
			f.Tags["bunny"] += "__,bind"
			if f.Nullable {
				f.Tags["bunny"] += ",null:" + f.Name
			}
		}
	}
	if _, ok := f.Tags["json"]; !ok {
		f.Tags["json"] = f.Name
		if f.Nullable {
			f.Tags["json"] += ",omitempty"
		}
	}
	return f.Tags.String()
}

func (f *Field) HasTag(tag string) bool {
	_, ok := f.Tags[tag]
	return ok
}

func (f *Field) IsStruct() bool {
	_, ok := f.Type.(*Struct)
	return ok
}

func (f *Field) TypeGo() TypeGo {
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
