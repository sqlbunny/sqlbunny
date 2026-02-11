package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/geo"
)

// PointZ is a nullable geo.PointZ.
type PointZ struct {
	PointZ geo.PointZ
	Valid bool
}

// NewPointZ creates a new PointZ
func NewPointZ(f geo.PointZ, valid bool) PointZ {
	return PointZ{
		PointZ: f,
		Valid: valid,
	}
}

// PointZFrom creates a new PointZ that will always be valid.
func PointZFrom(f geo.PointZ) PointZ {
	return NewPointZ(f, true)
}

// PointZFromPtr creates a new PointZ that be null if f is nil.
func PointZFromPtr(f *geo.PointZ) PointZ {
	if f == nil {
		return NewPointZ(geo.PointZ{}, false)
	}
	return NewPointZ(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *PointZ) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.PointZ = geo.PointZ{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.PointZ); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f PointZ) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.PointZ)
}

// SetValid changes this PointZ's value and also sets it to be non-null.
func (f *PointZ) SetValid(n geo.PointZ) {
	f.PointZ = n
	f.Valid = true
}

// Ptr returns a pointer to this PointZ's value, or a nil pointer if this PointZ is null.
func (f PointZ) Ptr() *geo.PointZ {
	if !f.Valid {
		return nil
	}
	return &f.PointZ
}

// IsZero returns true for invalid PointZs, for future omitempty support (Go 1.4?)
func (f PointZ) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *PointZ) Scan(value interface{}) error {
	if value == nil {
		f.PointZ, f.Valid = geo.PointZ{}, false
		return nil
	}
	f.Valid = true
	return f.PointZ.Scan(value)
}

// Value implements the driver Valuer interface.
func (f PointZ) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.PointZ.Value()
}
