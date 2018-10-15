package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// MultiPolygonZ is a nullable geo.MultiPolygonZ.
type MultiPolygonZ struct {
	MultiPolygonZ geo.MultiPolygonZ
	Valid bool
}

// NewMultiPolygonZ creates a new MultiPolygonZ
func NewMultiPolygonZ(f geo.MultiPolygonZ, valid bool) MultiPolygonZ {
	return MultiPolygonZ{
		MultiPolygonZ: f,
		Valid: valid,
	}
}

// MultiPolygonZFrom creates a new MultiPolygonZ that will always be valid.
func MultiPolygonZFrom(f geo.MultiPolygonZ) MultiPolygonZ {
	return NewMultiPolygonZ(f, true)
}

// MultiPolygonZFromPtr creates a new MultiPolygonZ that be null if f is nil.
func MultiPolygonZFromPtr(f *geo.MultiPolygonZ) MultiPolygonZ {
	if f == nil {
		return NewMultiPolygonZ(geo.MultiPolygonZ{}, false)
	}
	return NewMultiPolygonZ(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPolygonZ) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPolygonZ = geo.MultiPolygonZ{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPolygonZ); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPolygonZ) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPolygonZ)
}

// SetValid changes this MultiPolygonZ's value and also sets it to be non-null.
func (f *MultiPolygonZ) SetValid(n geo.MultiPolygonZ) {
	f.MultiPolygonZ = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPolygonZ's value, or a nil pointer if this MultiPolygonZ is null.
func (f MultiPolygonZ) Ptr() *geo.MultiPolygonZ {
	if !f.Valid {
		return nil
	}
	return &f.MultiPolygonZ
}

// IsZero returns true for invalid MultiPolygonZs, for future omitempty support (Go 1.4?)
func (f MultiPolygonZ) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPolygonZ) Scan(value interface{}) error {
	if value == nil {
		f.MultiPolygonZ, f.Valid = geo.MultiPolygonZ{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPolygonZ.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPolygonZ) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPolygonZ.Value()
}
