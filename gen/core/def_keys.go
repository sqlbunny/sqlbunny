package core

type modelPrimaryKey struct {
	names []string
}

func (modelPrimaryKey) isModelItem() {}

type fieldPrimaryKey func(...string) modelPrimaryKey

func (fieldPrimaryKey) isFieldItem() {}

var PrimaryKey fieldPrimaryKey = func(names ...string) modelPrimaryKey {
	return modelPrimaryKey{names: names}
}

type modelIndex struct {
	names []string
}

func (modelIndex) isModelItem() {}

type fieldIndex func(...string) modelIndex

func (fieldIndex) isFieldItem() {}

var Index fieldIndex = func(names ...string) modelIndex {
	return modelIndex{names: names}
}

type modelUnique struct {
	names []string
}

func (modelUnique) isModelItem() {}

type fieldUnique func(...string) modelUnique

func (fieldUnique) isFieldItem() {}

var Unique fieldUnique = func(names ...string) modelUnique {
	return modelUnique{names: names}
}

type modelForeignKey struct {
	foreignModelName   string
	columnNames        []string
	foreignColumnNames []string
}

func (modelForeignKey) isModelItem() {}

type fieldForeignKey struct {
	foreignModelName string
}

func (fieldForeignKey) isFieldItem() {}

func ForeignKey(foreignModelName string) fieldForeignKey {
	return fieldForeignKey{
		foreignModelName: foreignModelName,
	}
}

func ModelForeignKey(foreignModelName string, columnNames ...string) modelForeignKey {
	return modelForeignKey{
		foreignModelName: foreignModelName,
		columnNames:      columnNames,
	}
}
