package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// LineString is a nullable geo.LineString.
type LineString struct {
	LineString geo.LineString
	Valid bool
}

// NewLineString creates a new LineString
func NewLineString(f geo.LineString, valid bool) LineString {
	return LineString{
		LineString: f,
		Valid: valid,
	}
}

// LineStringFrom creates a new LineString that will always be valid.
func LineStringFrom(f geo.LineString) LineString {
	return NewLineString(f, true)
}

// LineStringFromPtr creates a new LineString that be null if f is nil.
func LineStringFromPtr(f *geo.LineString) LineString {
	if f == nil {
		return NewLineString(geo.LineString{}, false)
	}
	return NewLineString(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *LineString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.LineString = geo.LineString{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.LineString); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f LineString) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.LineString)
}

// SetValid changes this LineString's value and also sets it to be non-null.
func (f *LineString) SetValid(n geo.LineString) {
	f.LineString = n
	f.Valid = true
}

// Ptr returns a pointer to this LineString's value, or a nil pointer if this LineString is null.
func (f LineString) Ptr() *geo.LineString {
	if !f.Valid {
		return nil
	}
	return &f.LineString
}

// IsZero returns true for invalid LineStrings, for future omitempty support (Go 1.4?)
func (f LineString) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *LineString) Scan(value interface{}) error {
	if value == nil {
		f.LineString, f.Valid = geo.LineString{}, false
		return nil
	}
	f.Valid = true
	return f.LineString.Scan(value)
}

// Value implements the driver Valuer interface.
func (f LineString) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.LineString.Value()
}
