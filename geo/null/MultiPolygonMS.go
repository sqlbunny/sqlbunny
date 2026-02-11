package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// MultiPolygonMS is a nullable geo.MultiPolygonMS.
type MultiPolygonMS struct {
	MultiPolygonMS geo.MultiPolygonMS
	Valid bool
}

// NewMultiPolygonMS creates a new MultiPolygonMS
func NewMultiPolygonMS(f geo.MultiPolygonMS, valid bool) MultiPolygonMS {
	return MultiPolygonMS{
		MultiPolygonMS: f,
		Valid: valid,
	}
}

// MultiPolygonMSFrom creates a new MultiPolygonMS that will always be valid.
func MultiPolygonMSFrom(f geo.MultiPolygonMS) MultiPolygonMS {
	return NewMultiPolygonMS(f, true)
}

// MultiPolygonMSFromPtr creates a new MultiPolygonMS that be null if f is nil.
func MultiPolygonMSFromPtr(f *geo.MultiPolygonMS) MultiPolygonMS {
	if f == nil {
		return NewMultiPolygonMS(geo.MultiPolygonMS{}, false)
	}
	return NewMultiPolygonMS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPolygonMS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPolygonMS = geo.MultiPolygonMS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPolygonMS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPolygonMS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPolygonMS)
}

// SetValid changes this MultiPolygonMS's value and also sets it to be non-null.
func (f *MultiPolygonMS) SetValid(n geo.MultiPolygonMS) {
	f.MultiPolygonMS = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPolygonMS's value, or a nil pointer if this MultiPolygonMS is null.
func (f MultiPolygonMS) Ptr() *geo.MultiPolygonMS {
	if !f.Valid {
		return nil
	}
	return &f.MultiPolygonMS
}

// IsZero returns true for invalid MultiPolygonMSs, for future omitempty support (Go 1.4?)
func (f MultiPolygonMS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPolygonMS) Scan(value interface{}) error {
	if value == nil {
		f.MultiPolygonMS, f.Valid = geo.MultiPolygonMS{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPolygonMS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPolygonMS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPolygonMS.Value()
}
