package bunnyid

import "github.com/sqlbunny/sqlbunny/gen/core"
import "github.com/sqlbunny/sqlbunny/schema"

type ID struct {
	Prefix string
}

func (t ID) GetType(name string) schema.Type {
	return &IDType{
		Name:   name,
		Prefix: t.Prefix,
	}
}

func (t ID) ResolveTypes(v *core.Validation, st schema.Type, resolve func(name string, context string) schema.Type) {
}
