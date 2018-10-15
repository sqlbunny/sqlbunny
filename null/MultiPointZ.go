package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// MultiPointZ is a nullable geo.MultiPointZ.
type MultiPointZ struct {
	MultiPointZ geo.MultiPointZ
	Valid bool
}

// NewMultiPointZ creates a new MultiPointZ
func NewMultiPointZ(f geo.MultiPointZ, valid bool) MultiPointZ {
	return MultiPointZ{
		MultiPointZ: f,
		Valid: valid,
	}
}

// MultiPointZFrom creates a new MultiPointZ that will always be valid.
func MultiPointZFrom(f geo.MultiPointZ) MultiPointZ {
	return NewMultiPointZ(f, true)
}

// MultiPointZFromPtr creates a new MultiPointZ that be null if f is nil.
func MultiPointZFromPtr(f *geo.MultiPointZ) MultiPointZ {
	if f == nil {
		return NewMultiPointZ(geo.MultiPointZ{}, false)
	}
	return NewMultiPointZ(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPointZ) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPointZ = geo.MultiPointZ{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPointZ); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPointZ) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPointZ)
}

// SetValid changes this MultiPointZ's value and also sets it to be non-null.
func (f *MultiPointZ) SetValid(n geo.MultiPointZ) {
	f.MultiPointZ = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPointZ's value, or a nil pointer if this MultiPointZ is null.
func (f MultiPointZ) Ptr() *geo.MultiPointZ {
	if !f.Valid {
		return nil
	}
	return &f.MultiPointZ
}

// IsZero returns true for invalid MultiPointZs, for future omitempty support (Go 1.4?)
func (f MultiPointZ) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPointZ) Scan(value interface{}) error {
	if value == nil {
		f.MultiPointZ, f.Valid = geo.MultiPointZ{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPointZ.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPointZ) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPointZ.Value()
}
