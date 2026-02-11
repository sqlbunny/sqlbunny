package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/geo"
)

// LineStringM is a nullable geo.LineStringM.
type LineStringM struct {
	LineStringM geo.LineStringM
	Valid bool
}

// NewLineStringM creates a new LineStringM
func NewLineStringM(f geo.LineStringM, valid bool) LineStringM {
	return LineStringM{
		LineStringM: f,
		Valid: valid,
	}
}

// LineStringMFrom creates a new LineStringM that will always be valid.
func LineStringMFrom(f geo.LineStringM) LineStringM {
	return NewLineStringM(f, true)
}

// LineStringMFromPtr creates a new LineStringM that be null if f is nil.
func LineStringMFromPtr(f *geo.LineStringM) LineStringM {
	if f == nil {
		return NewLineStringM(geo.LineStringM{}, false)
	}
	return NewLineStringM(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *LineStringM) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.LineStringM = geo.LineStringM{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.LineStringM); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f LineStringM) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.LineStringM)
}

// SetValid changes this LineStringM's value and also sets it to be non-null.
func (f *LineStringM) SetValid(n geo.LineStringM) {
	f.LineStringM = n
	f.Valid = true
}

// Ptr returns a pointer to this LineStringM's value, or a nil pointer if this LineStringM is null.
func (f LineStringM) Ptr() *geo.LineStringM {
	if !f.Valid {
		return nil
	}
	return &f.LineStringM
}

// IsZero returns true for invalid LineStringMs, for future omitempty support (Go 1.4?)
func (f LineStringM) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *LineStringM) Scan(value interface{}) error {
	if value == nil {
		f.LineStringM, f.Valid = geo.LineStringM{}, false
		return nil
	}
	f.Valid = true
	return f.LineStringM.Scan(value)
}

// Value implements the driver Valuer interface.
func (f LineStringM) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.LineStringM.Value()
}
