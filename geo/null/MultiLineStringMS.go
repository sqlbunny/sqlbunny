package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/geo"
)

// MultiLineStringMS is a nullable geo.MultiLineStringMS.
type MultiLineStringMS struct {
	MultiLineStringMS geo.MultiLineStringMS
	Valid bool
}

// NewMultiLineStringMS creates a new MultiLineStringMS
func NewMultiLineStringMS(f geo.MultiLineStringMS, valid bool) MultiLineStringMS {
	return MultiLineStringMS{
		MultiLineStringMS: f,
		Valid: valid,
	}
}

// MultiLineStringMSFrom creates a new MultiLineStringMS that will always be valid.
func MultiLineStringMSFrom(f geo.MultiLineStringMS) MultiLineStringMS {
	return NewMultiLineStringMS(f, true)
}

// MultiLineStringMSFromPtr creates a new MultiLineStringMS that be null if f is nil.
func MultiLineStringMSFromPtr(f *geo.MultiLineStringMS) MultiLineStringMS {
	if f == nil {
		return NewMultiLineStringMS(geo.MultiLineStringMS{}, false)
	}
	return NewMultiLineStringMS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiLineStringMS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiLineStringMS = geo.MultiLineStringMS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiLineStringMS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiLineStringMS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiLineStringMS)
}

// SetValid changes this MultiLineStringMS's value and also sets it to be non-null.
func (f *MultiLineStringMS) SetValid(n geo.MultiLineStringMS) {
	f.MultiLineStringMS = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiLineStringMS's value, or a nil pointer if this MultiLineStringMS is null.
func (f MultiLineStringMS) Ptr() *geo.MultiLineStringMS {
	if !f.Valid {
		return nil
	}
	return &f.MultiLineStringMS
}

// IsZero returns true for invalid MultiLineStringMSs, for future omitempty support (Go 1.4?)
func (f MultiLineStringMS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiLineStringMS) Scan(value interface{}) error {
	if value == nil {
		f.MultiLineStringMS, f.Valid = geo.MultiLineStringMS{}, false
		return nil
	}
	f.Valid = true
	return f.MultiLineStringMS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiLineStringMS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiLineStringMS.Value()
}
