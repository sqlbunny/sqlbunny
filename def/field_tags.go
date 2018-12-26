package def

type FieldFlag interface {
	isFieldFlag()
}

type nullFlag struct{}

func (nullFlag) isFieldFlag() {}

var Null nullFlag

type primaryKeyFlag struct{}

func (primaryKeyFlag) isFieldFlag() {}

var PrimaryKey primaryKeyFlag

type indexFlag struct{}

func (indexFlag) isFieldFlag() {}

var Index indexFlag

type uniqueFlag struct{}

func (uniqueFlag) isFieldFlag() {}

var Unique uniqueFlag

type foreignKeyFlag struct {
	foreignModelName string
}

func (foreignKeyFlag) isFieldFlag() {}

func ForeignKey(foreignModelName string) foreignKeyFlag {
	return foreignKeyFlag{
		foreignModelName: foreignModelName,
	}
}

type tagFlag struct {
	key   string
	value string
}

func (tagFlag) isFieldFlag() {}

func Tag(key string, value string) tagFlag {
	return tagFlag{key: key, value: value}
}
