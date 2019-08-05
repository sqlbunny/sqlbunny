package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// MultiLineString is a nullable geo.MultiLineString.
type MultiLineString struct {
	MultiLineString geo.MultiLineString
	Valid bool
}

// NewMultiLineString creates a new MultiLineString
func NewMultiLineString(f geo.MultiLineString, valid bool) MultiLineString {
	return MultiLineString{
		MultiLineString: f,
		Valid: valid,
	}
}

// MultiLineStringFrom creates a new MultiLineString that will always be valid.
func MultiLineStringFrom(f geo.MultiLineString) MultiLineString {
	return NewMultiLineString(f, true)
}

// MultiLineStringFromPtr creates a new MultiLineString that be null if f is nil.
func MultiLineStringFromPtr(f *geo.MultiLineString) MultiLineString {
	if f == nil {
		return NewMultiLineString(geo.MultiLineString{}, false)
	}
	return NewMultiLineString(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiLineString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiLineString = geo.MultiLineString{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiLineString); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiLineString) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiLineString)
}

// SetValid changes this MultiLineString's value and also sets it to be non-null.
func (f *MultiLineString) SetValid(n geo.MultiLineString) {
	f.MultiLineString = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiLineString's value, or a nil pointer if this MultiLineString is null.
func (f MultiLineString) Ptr() *geo.MultiLineString {
	if !f.Valid {
		return nil
	}
	return &f.MultiLineString
}

// IsZero returns true for invalid MultiLineStrings, for future omitempty support (Go 1.4?)
func (f MultiLineString) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiLineString) Scan(value interface{}) error {
	if value == nil {
		f.MultiLineString, f.Valid = geo.MultiLineString{}, false
		return nil
	}
	f.Valid = true
	return f.MultiLineString.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiLineString) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiLineString.Value()
}
