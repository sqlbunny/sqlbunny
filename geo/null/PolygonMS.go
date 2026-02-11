package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// PolygonMS is a nullable geo.PolygonMS.
type PolygonMS struct {
	PolygonMS geo.PolygonMS
	Valid bool
}

// NewPolygonMS creates a new PolygonMS
func NewPolygonMS(f geo.PolygonMS, valid bool) PolygonMS {
	return PolygonMS{
		PolygonMS: f,
		Valid: valid,
	}
}

// PolygonMSFrom creates a new PolygonMS that will always be valid.
func PolygonMSFrom(f geo.PolygonMS) PolygonMS {
	return NewPolygonMS(f, true)
}

// PolygonMSFromPtr creates a new PolygonMS that be null if f is nil.
func PolygonMSFromPtr(f *geo.PolygonMS) PolygonMS {
	if f == nil {
		return NewPolygonMS(geo.PolygonMS{}, false)
	}
	return NewPolygonMS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *PolygonMS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.PolygonMS = geo.PolygonMS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.PolygonMS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f PolygonMS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.PolygonMS)
}

// SetValid changes this PolygonMS's value and also sets it to be non-null.
func (f *PolygonMS) SetValid(n geo.PolygonMS) {
	f.PolygonMS = n
	f.Valid = true
}

// Ptr returns a pointer to this PolygonMS's value, or a nil pointer if this PolygonMS is null.
func (f PolygonMS) Ptr() *geo.PolygonMS {
	if !f.Valid {
		return nil
	}
	return &f.PolygonMS
}

// IsZero returns true for invalid PolygonMSs, for future omitempty support (Go 1.4?)
func (f PolygonMS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *PolygonMS) Scan(value interface{}) error {
	if value == nil {
		f.PolygonMS, f.Valid = geo.PolygonMS{}, false
		return nil
	}
	f.Valid = true
	return f.PolygonMS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f PolygonMS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.PolygonMS.Value()
}
