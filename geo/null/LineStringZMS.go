package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/geo"
)

// LineStringZMS is a nullable geo.LineStringZMS.
type LineStringZMS struct {
	LineStringZMS geo.LineStringZMS
	Valid bool
}

// NewLineStringZMS creates a new LineStringZMS
func NewLineStringZMS(f geo.LineStringZMS, valid bool) LineStringZMS {
	return LineStringZMS{
		LineStringZMS: f,
		Valid: valid,
	}
}

// LineStringZMSFrom creates a new LineStringZMS that will always be valid.
func LineStringZMSFrom(f geo.LineStringZMS) LineStringZMS {
	return NewLineStringZMS(f, true)
}

// LineStringZMSFromPtr creates a new LineStringZMS that be null if f is nil.
func LineStringZMSFromPtr(f *geo.LineStringZMS) LineStringZMS {
	if f == nil {
		return NewLineStringZMS(geo.LineStringZMS{}, false)
	}
	return NewLineStringZMS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *LineStringZMS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.LineStringZMS = geo.LineStringZMS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.LineStringZMS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f LineStringZMS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.LineStringZMS)
}

// SetValid changes this LineStringZMS's value and also sets it to be non-null.
func (f *LineStringZMS) SetValid(n geo.LineStringZMS) {
	f.LineStringZMS = n
	f.Valid = true
}

// Ptr returns a pointer to this LineStringZMS's value, or a nil pointer if this LineStringZMS is null.
func (f LineStringZMS) Ptr() *geo.LineStringZMS {
	if !f.Valid {
		return nil
	}
	return &f.LineStringZMS
}

// IsZero returns true for invalid LineStringZMSs, for future omitempty support (Go 1.4?)
func (f LineStringZMS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *LineStringZMS) Scan(value interface{}) error {
	if value == nil {
		f.LineStringZMS, f.Valid = geo.LineStringZMS{}, false
		return nil
	}
	f.Valid = true
	return f.LineStringZMS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f LineStringZMS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.LineStringZMS.Value()
}
