package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// MultiLineStringZ is a nullable geo.MultiLineStringZ.
type MultiLineStringZ struct {
	MultiLineStringZ geo.MultiLineStringZ
	Valid bool
}

// NewMultiLineStringZ creates a new MultiLineStringZ
func NewMultiLineStringZ(f geo.MultiLineStringZ, valid bool) MultiLineStringZ {
	return MultiLineStringZ{
		MultiLineStringZ: f,
		Valid: valid,
	}
}

// MultiLineStringZFrom creates a new MultiLineStringZ that will always be valid.
func MultiLineStringZFrom(f geo.MultiLineStringZ) MultiLineStringZ {
	return NewMultiLineStringZ(f, true)
}

// MultiLineStringZFromPtr creates a new MultiLineStringZ that be null if f is nil.
func MultiLineStringZFromPtr(f *geo.MultiLineStringZ) MultiLineStringZ {
	if f == nil {
		return NewMultiLineStringZ(geo.MultiLineStringZ{}, false)
	}
	return NewMultiLineStringZ(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiLineStringZ) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiLineStringZ = geo.MultiLineStringZ{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiLineStringZ); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiLineStringZ) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiLineStringZ)
}

// SetValid changes this MultiLineStringZ's value and also sets it to be non-null.
func (f *MultiLineStringZ) SetValid(n geo.MultiLineStringZ) {
	f.MultiLineStringZ = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiLineStringZ's value, or a nil pointer if this MultiLineStringZ is null.
func (f MultiLineStringZ) Ptr() *geo.MultiLineStringZ {
	if !f.Valid {
		return nil
	}
	return &f.MultiLineStringZ
}

// IsZero returns true for invalid MultiLineStringZs, for future omitempty support (Go 1.4?)
func (f MultiLineStringZ) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiLineStringZ) Scan(value interface{}) error {
	if value == nil {
		f.MultiLineStringZ, f.Valid = geo.MultiLineStringZ{}, false
		return nil
	}
	f.Valid = true
	return f.MultiLineStringZ.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiLineStringZ) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiLineStringZ.Value()
}
