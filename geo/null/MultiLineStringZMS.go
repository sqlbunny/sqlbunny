package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/geo"
)

// MultiLineStringZMS is a nullable geo.MultiLineStringZMS.
type MultiLineStringZMS struct {
	MultiLineStringZMS geo.MultiLineStringZMS
	Valid bool
}

// NewMultiLineStringZMS creates a new MultiLineStringZMS
func NewMultiLineStringZMS(f geo.MultiLineStringZMS, valid bool) MultiLineStringZMS {
	return MultiLineStringZMS{
		MultiLineStringZMS: f,
		Valid: valid,
	}
}

// MultiLineStringZMSFrom creates a new MultiLineStringZMS that will always be valid.
func MultiLineStringZMSFrom(f geo.MultiLineStringZMS) MultiLineStringZMS {
	return NewMultiLineStringZMS(f, true)
}

// MultiLineStringZMSFromPtr creates a new MultiLineStringZMS that be null if f is nil.
func MultiLineStringZMSFromPtr(f *geo.MultiLineStringZMS) MultiLineStringZMS {
	if f == nil {
		return NewMultiLineStringZMS(geo.MultiLineStringZMS{}, false)
	}
	return NewMultiLineStringZMS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiLineStringZMS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiLineStringZMS = geo.MultiLineStringZMS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiLineStringZMS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiLineStringZMS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiLineStringZMS)
}

// SetValid changes this MultiLineStringZMS's value and also sets it to be non-null.
func (f *MultiLineStringZMS) SetValid(n geo.MultiLineStringZMS) {
	f.MultiLineStringZMS = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiLineStringZMS's value, or a nil pointer if this MultiLineStringZMS is null.
func (f MultiLineStringZMS) Ptr() *geo.MultiLineStringZMS {
	if !f.Valid {
		return nil
	}
	return &f.MultiLineStringZMS
}

// IsZero returns true for invalid MultiLineStringZMSs, for future omitempty support (Go 1.4?)
func (f MultiLineStringZMS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiLineStringZMS) Scan(value interface{}) error {
	if value == nil {
		f.MultiLineStringZMS, f.Valid = geo.MultiLineStringZMS{}, false
		return nil
	}
	f.Valid = true
	return f.MultiLineStringZMS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiLineStringZMS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiLineStringZMS.Value()
}
