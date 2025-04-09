package schema

import "strings"

type Type interface {
	GetName() string
	GoType() GoType

	GetExtension(key any) any
	SetExtension(key any, value any)
}

type BaseType interface {
	Type

	SQLType() SQLType
}

type NullableType interface {
	Type

	GoTypeNull() GoType
	GoTypeNullField() string
}

type GoType struct {
	Pkg  string
	Name string
}

type BaseTypeNotNullable struct {
	Name     string
	Go       GoType
	Postgres SQLType

	Extendable
}

type SQLType struct {
	Type      string
	ZeroValue string
}

func (t *BaseTypeNotNullable) GetName() string {
	return t.Name
}

func (t *BaseTypeNotNullable) GoType() GoType {
	return t.Go
}

func (t *BaseTypeNotNullable) SQLType() SQLType {
	return t.Postgres
}

type BaseTypeNullable struct {
	Name     string
	Go       GoType
	GoNull   GoType
	Postgres SQLType

	Extendable
}

func (t *BaseTypeNullable) GetName() string {
	return t.Name
}
func (t *BaseTypeNullable) GoType() GoType {
	return t.Go
}
func (t *BaseTypeNullable) GoTypeNull() GoType {
	return t.GoNull
}
func (t *BaseTypeNullable) GoTypeNullField() string {
	if strings.HasPrefix(t.GoNull.Name, "Null") {
		return t.GoNull.Name[4:]
	}
	return t.GoNull.Name
}
func (t *BaseTypeNullable) SQLType() SQLType {
	return t.Postgres
}

var _ BaseType = &BaseTypeNullable{}
var _ BaseType = &BaseTypeNotNullable{}
