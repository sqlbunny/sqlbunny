package schema

type Type interface {
	GetName() string
	TypeGo() TypeGo
}

type BaseType interface {
	GetName() string
	TypeGo() TypeGo
	TypeDB() string
}

type NullableType interface {
	GetName() string
	TypeGo() TypeGo
	TypeGoNull() TypeGo
	TypeGoNullField() string
}

type TypeGo struct {
	Pkg  string
	Name string
}

type BaseTypeNotNullable struct {
	Name     string
	Go       TypeGo
	Postgres string
}

func (t *BaseTypeNotNullable) GetName() string {
	return t.Name
}

func (t *BaseTypeNotNullable) TypeGo() TypeGo {
	return t.Go
}

func (t *BaseTypeNotNullable) TypeDB() string {
	return t.Postgres
}

type BaseTypeNullable struct {
	Name        string
	Go          TypeGo
	GoNull      TypeGo
	GoNullField string
	Postgres    string
}

func (t *BaseTypeNullable) GetName() string {
	return t.Name
}
func (t *BaseTypeNullable) TypeGo() TypeGo {
	return t.Go
}
func (t *BaseTypeNullable) TypeGoNull() TypeGo {
	return t.GoNull
}
func (t *BaseTypeNullable) TypeGoNullField() string {
	return t.GoNullField
}
func (t *BaseTypeNullable) TypeDB() string {
	return t.Postgres
}
