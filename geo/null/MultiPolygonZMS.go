package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// MultiPolygonZMS is a nullable geo.MultiPolygonZMS.
type MultiPolygonZMS struct {
	MultiPolygonZMS geo.MultiPolygonZMS
	Valid bool
}

// NewMultiPolygonZMS creates a new MultiPolygonZMS
func NewMultiPolygonZMS(f geo.MultiPolygonZMS, valid bool) MultiPolygonZMS {
	return MultiPolygonZMS{
		MultiPolygonZMS: f,
		Valid: valid,
	}
}

// MultiPolygonZMSFrom creates a new MultiPolygonZMS that will always be valid.
func MultiPolygonZMSFrom(f geo.MultiPolygonZMS) MultiPolygonZMS {
	return NewMultiPolygonZMS(f, true)
}

// MultiPolygonZMSFromPtr creates a new MultiPolygonZMS that be null if f is nil.
func MultiPolygonZMSFromPtr(f *geo.MultiPolygonZMS) MultiPolygonZMS {
	if f == nil {
		return NewMultiPolygonZMS(geo.MultiPolygonZMS{}, false)
	}
	return NewMultiPolygonZMS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPolygonZMS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPolygonZMS = geo.MultiPolygonZMS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPolygonZMS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPolygonZMS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPolygonZMS)
}

// SetValid changes this MultiPolygonZMS's value and also sets it to be non-null.
func (f *MultiPolygonZMS) SetValid(n geo.MultiPolygonZMS) {
	f.MultiPolygonZMS = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPolygonZMS's value, or a nil pointer if this MultiPolygonZMS is null.
func (f MultiPolygonZMS) Ptr() *geo.MultiPolygonZMS {
	if !f.Valid {
		return nil
	}
	return &f.MultiPolygonZMS
}

// IsZero returns true for invalid MultiPolygonZMSs, for future omitempty support (Go 1.4?)
func (f MultiPolygonZMS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPolygonZMS) Scan(value interface{}) error {
	if value == nil {
		f.MultiPolygonZMS, f.Valid = geo.MultiPolygonZMS{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPolygonZMS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPolygonZMS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPolygonZMS.Value()
}
