package schema

type Extendable struct {
	values map[any]any
}

func (e *Extendable) GetExtension(key any) any {
	if e.values == nil {
		return nil
	}
	return e.values[key]
}

func (e *Extendable) SetExtension(key any, value any) {
	if e.values == nil {
		e.values = make(map[any]any)
	}
	e.values[key] = value
}
