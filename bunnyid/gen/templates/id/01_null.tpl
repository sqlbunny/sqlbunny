{{- $modelName := .IDType.Name | titleCase -}}

// Null{{$modelName}} is an nullable {{$modelName}}.
type Null{{$modelName}} struct {
	ID    {{$modelName}}
	Valid bool
}

// NewNull{{$modelName}} creates a new Null{{$modelName}}
func NewNull{{$modelName}}(i {{$modelName}}, valid bool) Null{{$modelName}} {
	return Null{{$modelName}}{
		ID:    i,
		Valid: valid,
	}
}

// Null{{$modelName}}From creates a new ID that will always be valid.
func Null{{$modelName}}From(i {{$modelName}}) Null{{$modelName}} {
	return NewNull{{$modelName}}(i, true)
}

// Null{{$modelName}}FromPtr creates a new ID that be null if i is nil.
func Null{{$modelName}}FromPtr(i *{{$modelName}}) Null{{$modelName}} {
	if i == nil {
		return NewNull{{$modelName}}({{$modelName}}{}, false)
	}
	return NewNull{{$modelName}}(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *Null{{$modelName}}) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		u.ID = {{$modelName}}{}
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.ID); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *Null{{$modelName}}) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := {{$modelName}}FromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.ID = res
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (u Null{{$modelName}}) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return u.ID.MarshalJSON()
}

// MarshalText implements encoding.TextMarshaler.
func (u Null{{$modelName}}) MarshalText() ([]byte, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.MarshalText()
}

// SetValid changes this ID's value and also sets it to be non-null.
func (u *Null{{$modelName}}) SetValid(n {{$modelName}}) {
	u.ID = n
	u.Valid = true
}

// Ptr returns a pointer to this ID's value, or a nil pointer if this ID is null.
func (u Null{{$modelName}}) Ptr() *{{$modelName}} {
	if !u.Valid {
		return nil
	}
	return &u.ID
}

// IsZero returns true for invalid ID's, for future omitempty support (Go 1.4?)
func (u Null{{$modelName}}) IsZero() bool {
	return !u.Valid
}

// Scan implements the Scanner interface.
func (u *Null{{$modelName}}) Scan(value interface{}) error {
	if value == nil {
		u.ID, u.Valid = {{$modelName}}{}, false
		return nil
	}
	u.Valid = true
	return u.ID.Scan(value)
}

// Value implements the driver Valuer interface.
func (u Null{{$modelName}}) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.Value()
}

func (u Null{{$modelName}}) String() string {
	if !u.Valid {
		return "<null {{$modelName}}>"
	}
	return u.ID.String()
}
