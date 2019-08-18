package schema

type Extendable struct {
	values map[interface{}]interface{}
}

func (e *Extendable) GetExtension(key interface{}) interface{} {
	if e.values == nil {
		return nil
	}
	return e.values[key]
}

func (e *Extendable) SetExtension(key interface{}, value interface{}) {
	if e.values == nil {
		e.values = make(map[interface{}]interface{})
	}
	e.values[key] = value
}
