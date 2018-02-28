{{- $modelName := .IDType.Name | titleCase -}}


// ID is an nullable ID.
type Null{{$modelName}}ID struct {
	ID    {{$modelName}}ID
	Valid bool
}

// NewNull{{$modelName}}ID creates a new Null{{$modelName}}ID
func NewNull{{$modelName}}ID(i {{$modelName}}ID, valid bool) Null{{$modelName}}ID {
	return Null{{$modelName}}ID{
		ID:    i,
		Valid: valid,
	}
}

// Null{{$modelName}}IDFrom creates a new ID that will always be valid.
func Null{{$modelName}}IDFrom(i {{$modelName}}ID) Null{{$modelName}}ID {
	return NewNull{{$modelName}}ID(i, true)
}

// Null{{$modelName}}IDFromPtr creates a new ID that be null if i is nil.
func Null{{$modelName}}IDFromPtr(i *{{$modelName}}ID) Null{{$modelName}}ID {
	if i == nil {
		return NewNull{{$modelName}}ID({{$modelName}}ID{}, false)
	}
	return NewNull{{$modelName}}ID(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *Null{{$modelName}}ID) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, boil.IDNullBytes) {
		u.ID = {{$modelName}}ID{}
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
func (u *Null{{$modelName}}ID) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := {{$modelName}}IDFromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.ID = res
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (u Null{{$modelName}}ID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return boil.IDNullBytes, nil
	}
	return u.ID.MarshalJSON()
}

// MarshalText implements encoding.TextMarshaler.
func (u Null{{$modelName}}ID) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return u.ID.MarshalText()
}

// SetValid changes this ID's value and also sets it to be non-null.
func (u *Null{{$modelName}}ID) SetValid(n {{$modelName}}ID) {
	u.ID = n
	u.Valid = true
}

// Ptr returns a pointer to this ID's value, or a nil pointer if this ID is null.
func (u Null{{$modelName}}ID) Ptr() *{{$modelName}}ID {
	if !u.Valid {
		return nil
	}
	return &u.ID
}

// IsZero returns true for invalid ID's, for future omitempty support (Go 1.4?)
func (u Null{{$modelName}}ID) IsZero() bool {
	return !u.Valid
}

// Scan implements the Scanner interface.
func (u *Null{{$modelName}}ID) Scan(value interface{}) error {
	if value == nil {
		u.ID, u.Valid = {{$modelName}}ID{}, false
		return nil
	}
	u.Valid = true
	return u.ID.Scan(value)
}

// Value implements the driver Valuer interface.
func (u Null{{$modelName}}ID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.Value()
}
