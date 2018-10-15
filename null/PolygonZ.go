package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// PolygonZ is a nullable geo.PolygonZ.
type PolygonZ struct {
	PolygonZ geo.PolygonZ
	Valid bool
}

// NewPolygonZ creates a new PolygonZ
func NewPolygonZ(f geo.PolygonZ, valid bool) PolygonZ {
	return PolygonZ{
		PolygonZ: f,
		Valid: valid,
	}
}

// PolygonZFrom creates a new PolygonZ that will always be valid.
func PolygonZFrom(f geo.PolygonZ) PolygonZ {
	return NewPolygonZ(f, true)
}

// PolygonZFromPtr creates a new PolygonZ that be null if f is nil.
func PolygonZFromPtr(f *geo.PolygonZ) PolygonZ {
	if f == nil {
		return NewPolygonZ(geo.PolygonZ{}, false)
	}
	return NewPolygonZ(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *PolygonZ) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.PolygonZ = geo.PolygonZ{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.PolygonZ); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f PolygonZ) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.PolygonZ)
}

// SetValid changes this PolygonZ's value and also sets it to be non-null.
func (f *PolygonZ) SetValid(n geo.PolygonZ) {
	f.PolygonZ = n
	f.Valid = true
}

// Ptr returns a pointer to this PolygonZ's value, or a nil pointer if this PolygonZ is null.
func (f PolygonZ) Ptr() *geo.PolygonZ {
	if !f.Valid {
		return nil
	}
	return &f.PolygonZ
}

// IsZero returns true for invalid PolygonZs, for future omitempty support (Go 1.4?)
func (f PolygonZ) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *PolygonZ) Scan(value interface{}) error {
	if value == nil {
		f.PolygonZ, f.Valid = geo.PolygonZ{}, false
		return nil
	}
	f.Valid = true
	return f.PolygonZ.Scan(value)
}

// Value implements the driver Valuer interface.
func (f PolygonZ) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.PolygonZ.Value()
}
