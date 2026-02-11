package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// PolygonZMS is a nullable geo.PolygonZMS.
type PolygonZMS struct {
	PolygonZMS geo.PolygonZMS
	Valid bool
}

// NewPolygonZMS creates a new PolygonZMS
func NewPolygonZMS(f geo.PolygonZMS, valid bool) PolygonZMS {
	return PolygonZMS{
		PolygonZMS: f,
		Valid: valid,
	}
}

// PolygonZMSFrom creates a new PolygonZMS that will always be valid.
func PolygonZMSFrom(f geo.PolygonZMS) PolygonZMS {
	return NewPolygonZMS(f, true)
}

// PolygonZMSFromPtr creates a new PolygonZMS that be null if f is nil.
func PolygonZMSFromPtr(f *geo.PolygonZMS) PolygonZMS {
	if f == nil {
		return NewPolygonZMS(geo.PolygonZMS{}, false)
	}
	return NewPolygonZMS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *PolygonZMS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.PolygonZMS = geo.PolygonZMS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.PolygonZMS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f PolygonZMS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.PolygonZMS)
}

// SetValid changes this PolygonZMS's value and also sets it to be non-null.
func (f *PolygonZMS) SetValid(n geo.PolygonZMS) {
	f.PolygonZMS = n
	f.Valid = true
}

// Ptr returns a pointer to this PolygonZMS's value, or a nil pointer if this PolygonZMS is null.
func (f PolygonZMS) Ptr() *geo.PolygonZMS {
	if !f.Valid {
		return nil
	}
	return &f.PolygonZMS
}

// IsZero returns true for invalid PolygonZMSs, for future omitempty support (Go 1.4?)
func (f PolygonZMS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *PolygonZMS) Scan(value interface{}) error {
	if value == nil {
		f.PolygonZMS, f.Valid = geo.PolygonZMS{}, false
		return nil
	}
	f.Valid = true
	return f.PolygonZMS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f PolygonZMS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.PolygonZMS.Value()
}
