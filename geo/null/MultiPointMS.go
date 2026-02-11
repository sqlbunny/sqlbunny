package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/geo"
)

// MultiPointMS is a nullable geo.MultiPointMS.
type MultiPointMS struct {
	MultiPointMS geo.MultiPointMS
	Valid bool
}

// NewMultiPointMS creates a new MultiPointMS
func NewMultiPointMS(f geo.MultiPointMS, valid bool) MultiPointMS {
	return MultiPointMS{
		MultiPointMS: f,
		Valid: valid,
	}
}

// MultiPointMSFrom creates a new MultiPointMS that will always be valid.
func MultiPointMSFrom(f geo.MultiPointMS) MultiPointMS {
	return NewMultiPointMS(f, true)
}

// MultiPointMSFromPtr creates a new MultiPointMS that be null if f is nil.
func MultiPointMSFromPtr(f *geo.MultiPointMS) MultiPointMS {
	if f == nil {
		return NewMultiPointMS(geo.MultiPointMS{}, false)
	}
	return NewMultiPointMS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPointMS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPointMS = geo.MultiPointMS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPointMS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPointMS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPointMS)
}

// SetValid changes this MultiPointMS's value and also sets it to be non-null.
func (f *MultiPointMS) SetValid(n geo.MultiPointMS) {
	f.MultiPointMS = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPointMS's value, or a nil pointer if this MultiPointMS is null.
func (f MultiPointMS) Ptr() *geo.MultiPointMS {
	if !f.Valid {
		return nil
	}
	return &f.MultiPointMS
}

// IsZero returns true for invalid MultiPointMSs, for future omitempty support (Go 1.4?)
func (f MultiPointMS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPointMS) Scan(value interface{}) error {
	if value == nil {
		f.MultiPointMS, f.Valid = geo.MultiPointMS{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPointMS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPointMS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPointMS.Value()
}
