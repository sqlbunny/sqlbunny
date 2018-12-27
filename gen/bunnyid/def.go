package bunnyid

import "github.com/kernelpayments/sqlbunny/def"
import "github.com/kernelpayments/sqlbunny/schema"

type ID struct {
	Prefix string
}

func (t ID) GetType(name string) schema.Type {
	return &IDType{
		Name:   name,
		Prefix: t.Prefix,
	}
}

func (t ID) ResolveTypes(v *def.Validation, st schema.Type, resolve func(name string, context string) schema.Type) {
}
