package core

type FieldItem interface {
	isFieldItem()
}

type fieldNull struct{}

func (fieldNull) isFieldItem() {}

var Null fieldNull

type fieldTag struct {
	key   string
	value string
}

func (fieldTag) isFieldItem() {}

func Tag(key string, value string) fieldTag {
	return fieldTag{key: key, value: value}
}

type field struct {
	name     string
	typeName string
	flags    []FieldItem
}

func (field) isModelItem() {}

func Field(name string, typeName string, flags ...FieldItem) field {
	return field{
		name:     name,
		typeName: typeName,
		flags:    flags,
	}
}

type ModelItem interface {
	isModelItem()
}

type model struct {
	name  string
	items []ModelItem
}

func (model) IsConfigItem() {}

func expandItems(items []ModelItem) []ModelItem {
	res := items
	for _, i := range items {
		if i, ok := i.(field); ok {
			for _, f := range i.flags {
				switch f := f.(type) {
				case fieldPrimaryKey:
					res = append(res, modelPrimaryKey{names: []string{i.name}})
				case fieldUnique:
					res = append(res, modelUnique{names: []string{i.name}})
				case fieldIndex:
					res = append(res, modelIndex{names: []string{i.name}})
				case fieldForeignKey:
					res = append(res, modelForeignKey{
						columnName:       i.name,
						foreignModelName: f.foreignModelName,
					})
				}
			}
		}
	}
	return res
}

func Model(name string, items ...ModelItem) model {
	return model{
		name:  name,
		items: expandItems(items),
	}
}
