package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// LineStringZ is a nullable geo.LineStringZ.
type LineStringZ struct {
	LineStringZ geo.LineStringZ
	Valid bool
}

// NewLineStringZ creates a new LineStringZ
func NewLineStringZ(f geo.LineStringZ, valid bool) LineStringZ {
	return LineStringZ{
		LineStringZ: f,
		Valid: valid,
	}
}

// LineStringZFrom creates a new LineStringZ that will always be valid.
func LineStringZFrom(f geo.LineStringZ) LineStringZ {
	return NewLineStringZ(f, true)
}

// LineStringZFromPtr creates a new LineStringZ that be null if f is nil.
func LineStringZFromPtr(f *geo.LineStringZ) LineStringZ {
	if f == nil {
		return NewLineStringZ(geo.LineStringZ{}, false)
	}
	return NewLineStringZ(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *LineStringZ) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.LineStringZ = geo.LineStringZ{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.LineStringZ); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f LineStringZ) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.LineStringZ)
}

// SetValid changes this LineStringZ's value and also sets it to be non-null.
func (f *LineStringZ) SetValid(n geo.LineStringZ) {
	f.LineStringZ = n
	f.Valid = true
}

// Ptr returns a pointer to this LineStringZ's value, or a nil pointer if this LineStringZ is null.
func (f LineStringZ) Ptr() *geo.LineStringZ {
	if !f.Valid {
		return nil
	}
	return &f.LineStringZ
}

// IsZero returns true for invalid LineStringZs, for future omitempty support (Go 1.4?)
func (f LineStringZ) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *LineStringZ) Scan(value interface{}) error {
	if value == nil {
		f.LineStringZ, f.Valid = geo.LineStringZ{}, false
		return nil
	}
	f.Valid = true
	return f.LineStringZ.Scan(value)
}

// Value implements the driver Valuer interface.
func (f LineStringZ) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.LineStringZ.Value()
}
