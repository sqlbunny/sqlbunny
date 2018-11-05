package def

type field struct {
	name     string
	typeName string
	flags    []FieldFlag
}

func (field) isModelItem() {}

func Field(name string, typeName string, flags ...FieldFlag) field {
	return field{
		name:     name,
		typeName: typeName,
		flags:    flags,
	}
}

type ModelItem interface {
	isModelItem()
}

type primaryKey struct {
	names []string
}

func (primaryKey) isModelItem() {}

func ModelPrimaryKey(names ...string) primaryKey {
	return primaryKey{names: names}
}

type index struct {
	names []string
}

func (index) isModelItem() {}

func ModelIndex(names ...string) index {
	return index{names: names}
}

type unique struct {
	names []string
}

func (unique) isModelItem() {}

func ModelUnique(names ...string) unique {
	return unique{names: names}
}

type foreignKey struct {
	columnName       string
	foreignModelName string
}

func (foreignKey) isModelItem() {}

func ModelForeignKey(columnName, foreignModelName string) foreignKey {
	return foreignKey{
		columnName:       columnName,
		foreignModelName: foreignModelName,
	}
}

type model struct {
	name  string
	items []ModelItem
}

var models []model

func expandItems(items []ModelItem) []ModelItem {
	res := items
	for _, i := range items {
		if i, ok := i.(field); ok {
			for _, f := range i.flags {
				switch f := f.(type) {
				case primaryKeyFlag:
					res = append(res, primaryKey{names: []string{i.name}})
				case uniqueFlag:
					res = append(res, unique{names: []string{i.name}})
				case indexFlag:
					res = append(res, index{names: []string{i.name}})
				case foreignKeyFlag:
					res = append(res, foreignKey{
						columnName:       i.name,
						foreignModelName: f.foreignModelName,
					})
				}
			}
		}
	}
	return res
}

func Model(name string, items ...ModelItem) {
	models = append(models, model{
		name:  name,
		items: expandItems(items),
	})
}
