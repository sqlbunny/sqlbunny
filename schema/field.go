package schema

// Field holds information about a database field.
// Types are Go types, converted by TranslateFieldType.
type Field struct {
	Name string

	Type Type
	Tags Tags

	Nullable bool

	// Temp vars for use while parsing
	typeName   string
	index      bool
	unique     bool
	primaryKey bool
	foreignKey string
}

func (f *Field) GenerateTags() string {
	if _, ok := f.Tags["boil"]; !ok {
		f.Tags["boil"] = f.Name
		if f.IsStruct() {
			f.Tags["boil"] += ","
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

//generateTags $dot.Tags $field.Name}}boil:"{{$field.Name}}" json:"{{$field.Name}}{{if $field.Nullable}},omitempty{{end}}" toml:"{{$field.Name}}" yaml:"{{$field.Name}}{{if $field.Nullable}},omitempty{{end}}" {{$field.Tag}}

func (f *Field) HasTag(tag string) bool {
	_, ok := f.Tags[tag]
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
