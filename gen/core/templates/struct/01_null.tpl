{{- $dot := . -}}
{{- $modelName := .Struct.Name | titleCase -}}

type Null{{$modelName}} struct {
	{{$modelName}} {{$modelName}}
	Valid    bool
}

func NewNull{{$modelName}}(s {{$modelName}}, valid bool) Null{{$modelName}} {
	return Null{{$modelName}}{
		{{$modelName}}: s,
		Valid:    valid,
	}
}

func Null{{$modelName}}From(s {{$modelName}}) Null{{$modelName}} {
	return NewNull{{$modelName}}(s, true)
}

func Null{{$modelName}}FromPtr(s *{{$modelName}}) Null{{$modelName}} {
	if s == nil {
		return NewNull{{$modelName}}({{$modelName}}{}, false)
	}
	return NewNull{{$modelName}}(*s, true)
}

func (u *Null{{$modelName}}) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		u.{{$modelName}} = {{$modelName}}{}
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.{{$modelName}}); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

func (u Null{{$modelName}}) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return json.Marshal(u.{{$modelName}})
}

func (u *Null{{$modelName}}) SetValid(n {{$modelName}}) {
	u.{{$modelName}} = n
	u.Valid = true
}

func (u Null{{$modelName}}) Ptr() *{{$modelName}} {
	if !u.Valid {
		return nil
	}
	return &u.{{$modelName}}
}

func (u Null{{$modelName}}) IsZero() bool {
	return !u.Valid
}
