package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// LineStringMS is a nullable geo.LineStringMS.
type LineStringMS struct {
	LineStringMS geo.LineStringMS
	Valid bool
}

// NewLineStringMS creates a new LineStringMS
func NewLineStringMS(f geo.LineStringMS, valid bool) LineStringMS {
	return LineStringMS{
		LineStringMS: f,
		Valid: valid,
	}
}

// LineStringMSFrom creates a new LineStringMS that will always be valid.
func LineStringMSFrom(f geo.LineStringMS) LineStringMS {
	return NewLineStringMS(f, true)
}

// LineStringMSFromPtr creates a new LineStringMS that be null if f is nil.
func LineStringMSFromPtr(f *geo.LineStringMS) LineStringMS {
	if f == nil {
		return NewLineStringMS(geo.LineStringMS{}, false)
	}
	return NewLineStringMS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *LineStringMS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.LineStringMS = geo.LineStringMS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.LineStringMS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f LineStringMS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.LineStringMS)
}

// SetValid changes this LineStringMS's value and also sets it to be non-null.
func (f *LineStringMS) SetValid(n geo.LineStringMS) {
	f.LineStringMS = n
	f.Valid = true
}

// Ptr returns a pointer to this LineStringMS's value, or a nil pointer if this LineStringMS is null.
func (f LineStringMS) Ptr() *geo.LineStringMS {
	if !f.Valid {
		return nil
	}
	return &f.LineStringMS
}

// IsZero returns true for invalid LineStringMSs, for future omitempty support (Go 1.4?)
func (f LineStringMS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *LineStringMS) Scan(value interface{}) error {
	if value == nil {
		f.LineStringMS, f.Valid = geo.LineStringMS{}, false
		return nil
	}
	f.Valid = true
	return f.LineStringMS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f LineStringMS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.LineStringMS.Value()
}
