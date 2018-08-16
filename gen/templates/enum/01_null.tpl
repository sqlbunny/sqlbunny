{{- $dot := . -}}
{{- $enumName := .Enum.Name | titleCase -}}
{{- $enumNameCamel := .Enum.Name | camelCase -}}

// Null{{$enumName}} is a nullable {{$enumName}}.
type Null{{$enumName}} struct {
	{{$enumName}}    {{$enumName}}
	Valid bool
}

// NewNull{{$enumName}} creates a new Null{{$enumName}}
func NewNull{{$enumName}}(i {{$enumName}}, valid bool) Null{{$enumName}} {
	return Null{{$enumName}}{
		{{$enumName}}:    i,
		Valid: valid,
	}
}

// Null{{$enumName}}From creates a new Null{{$enumName}} that will always be valid.
func Null{{$enumName}}From(i {{$enumName}}) Null{{$enumName}} {
	return NewNull{{$enumName}}(i, true)
}

// Null{{$enumName}}FromPtr creates a new Null{{$enumName}} that be null if i is nil.
func Null{{$enumName}}FromPtr(i *{{$enumName}}) Null{{$enumName}} {
	if i == nil {
        var z {{$enumName}}
		return NewNull{{$enumName}}(z, false)
	}
	return NewNull{{$enumName}}(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *Null{{$enumName}}) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, boil.IDNullBytes) {
        var z {{$enumName}}
		u.{{$enumName}} = z
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.{{$enumName}}); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *Null{{$enumName}}) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := {{$enumName}}FromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.{{$enumName}} = res
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (u Null{{$enumName}}) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return boil.IDNullBytes, nil
	}
	return json.Marshal(u.{{$enumName}})
}

// MarshalText implements encoding.TextMarshaler.
func (u Null{{$enumName}}) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return u.{{$enumName}}.MarshalText()
}

// SetValid changes this {{$enumName}}'s value and also sets it to be non-null.
func (u *Null{{$enumName}}) SetValid(n {{$enumName}}) {
	u.{{$enumName}} = n
	u.Valid = true
}

// Ptr returns a pointer to this {{$enumName}}'s value, or a nil pointer if this {{$enumName}} is null.
func (u Null{{$enumName}}) Ptr() *{{$enumName}} {
	if !u.Valid {
		return nil
	}
	return &u.{{$enumName}}
}

// IsZero returns true for invalid {{$enumName}}'s, for future omitempty support (Go 1.4?)
func (u Null{{$enumName}}) IsZero() bool {
	return !u.Valid
}

// Scan implements the Scanner interface.
func (u *Null{{$enumName}}) Scan(value interface{}) error {
	if value == nil {
        var z {{$enumName}}
		u.{{$enumName}}, u.Valid = z, false
		return nil
	}
	u.Valid = true
    return convert.ConvertAssign(&u.{{$enumName}}, value)
}

// Value implements the driver Valuer interface.
func (u Null{{$enumName}}) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.{{$enumName}}, nil
}

func (u Null{{$enumName}}) String() string {
	if !u.Valid {
		return "<null {{$enumName}}>"
	}
	return u.{{$enumName}}.String()
}
