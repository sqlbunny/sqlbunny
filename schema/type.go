package schema

type Type interface {
	GetName() string
	TypeGo() string
}

type BaseType interface {
	GetName() string
	TypeGo() string
	TypeDB() string
}

type NullableType interface {
	GetName() string
	TypeGo() string
	TypeGoNull() string
	TypeGoNullField() string
}

type TypeWithImports interface {
	GetImports() []string
}

type BaseTypeNotNullable struct {
	Name      string
	Go        string
	Postgres  string
	GoImports []string
}

func (t *BaseTypeNotNullable) GetName() string {
	return t.Name
}

func (t *BaseTypeNotNullable) TypeGo() string {
	return t.Go
}

func (t *BaseTypeNotNullable) TypeDB() string {
	return t.Postgres
}

func (t *BaseTypeNotNullable) GetImports() []string {
	return t.GoImports
}

type BaseTypeNullable struct {
	Name        string
	Go          string
	GoNull      string
	GoNullField string
	Postgres    string
	GoImports   []string
}

func (t *BaseTypeNullable) GetName() string {
	return t.Name
}
func (t *BaseTypeNullable) TypeGo() string {
	return t.Go
}
func (t *BaseTypeNullable) TypeGoNull() string {
	return t.GoNull
}
func (t *BaseTypeNullable) TypeGoNullField() string {
	return t.GoNullField
}
func (t *BaseTypeNullable) TypeDB() string {
	return t.Postgres
}
func (t *BaseTypeNullable) GetImports() []string {
	return t.GoImports
}
