package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// MultiPolygon is a nullable geo.MultiPolygon.
type MultiPolygon struct {
	MultiPolygon geo.MultiPolygon
	Valid bool
}

// NewMultiPolygon creates a new MultiPolygon
func NewMultiPolygon(f geo.MultiPolygon, valid bool) MultiPolygon {
	return MultiPolygon{
		MultiPolygon: f,
		Valid: valid,
	}
}

// MultiPolygonFrom creates a new MultiPolygon that will always be valid.
func MultiPolygonFrom(f geo.MultiPolygon) MultiPolygon {
	return NewMultiPolygon(f, true)
}

// MultiPolygonFromPtr creates a new MultiPolygon that be null if f is nil.
func MultiPolygonFromPtr(f *geo.MultiPolygon) MultiPolygon {
	if f == nil {
		return NewMultiPolygon(geo.MultiPolygon{}, false)
	}
	return NewMultiPolygon(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPolygon) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPolygon = geo.MultiPolygon{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPolygon); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPolygon) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPolygon)
}

// SetValid changes this MultiPolygon's value and also sets it to be non-null.
func (f *MultiPolygon) SetValid(n geo.MultiPolygon) {
	f.MultiPolygon = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPolygon's value, or a nil pointer if this MultiPolygon is null.
func (f MultiPolygon) Ptr() *geo.MultiPolygon {
	if !f.Valid {
		return nil
	}
	return &f.MultiPolygon
}

// IsZero returns true for invalid MultiPolygons, for future omitempty support (Go 1.4?)
func (f MultiPolygon) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPolygon) Scan(value interface{}) error {
	if value == nil {
		f.MultiPolygon, f.Valid = geo.MultiPolygon{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPolygon.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPolygon) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPolygon.Value()
}
