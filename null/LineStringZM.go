package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// LineStringZM is a nullable geo.LineStringZM.
type LineStringZM struct {
	LineStringZM geo.LineStringZM
	Valid bool
}

// NewLineStringZM creates a new LineStringZM
func NewLineStringZM(f geo.LineStringZM, valid bool) LineStringZM {
	return LineStringZM{
		LineStringZM: f,
		Valid: valid,
	}
}

// LineStringZMFrom creates a new LineStringZM that will always be valid.
func LineStringZMFrom(f geo.LineStringZM) LineStringZM {
	return NewLineStringZM(f, true)
}

// LineStringZMFromPtr creates a new LineStringZM that be null if f is nil.
func LineStringZMFromPtr(f *geo.LineStringZM) LineStringZM {
	if f == nil {
		return NewLineStringZM(geo.LineStringZM{}, false)
	}
	return NewLineStringZM(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *LineStringZM) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.LineStringZM = geo.LineStringZM{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.LineStringZM); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f LineStringZM) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.LineStringZM)
}

// SetValid changes this LineStringZM's value and also sets it to be non-null.
func (f *LineStringZM) SetValid(n geo.LineStringZM) {
	f.LineStringZM = n
	f.Valid = true
}

// Ptr returns a pointer to this LineStringZM's value, or a nil pointer if this LineStringZM is null.
func (f LineStringZM) Ptr() *geo.LineStringZM {
	if !f.Valid {
		return nil
	}
	return &f.LineStringZM
}

// IsZero returns true for invalid LineStringZMs, for future omitempty support (Go 1.4?)
func (f LineStringZM) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *LineStringZM) Scan(value interface{}) error {
	if value == nil {
		f.LineStringZM, f.Valid = geo.LineStringZM{}, false
		return nil
	}
	f.Valid = true
	return f.LineStringZM.Scan(value)
}

// Value implements the driver Valuer interface.
func (f LineStringZM) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.LineStringZM.Value()
}
