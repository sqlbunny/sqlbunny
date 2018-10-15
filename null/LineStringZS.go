package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// LineStringZS is a nullable geo.LineStringZS.
type LineStringZS struct {
	LineStringZS geo.LineStringZS
	Valid bool
}

// NewLineStringZS creates a new LineStringZS
func NewLineStringZS(f geo.LineStringZS, valid bool) LineStringZS {
	return LineStringZS{
		LineStringZS: f,
		Valid: valid,
	}
}

// LineStringZSFrom creates a new LineStringZS that will always be valid.
func LineStringZSFrom(f geo.LineStringZS) LineStringZS {
	return NewLineStringZS(f, true)
}

// LineStringZSFromPtr creates a new LineStringZS that be null if f is nil.
func LineStringZSFromPtr(f *geo.LineStringZS) LineStringZS {
	if f == nil {
		return NewLineStringZS(geo.LineStringZS{}, false)
	}
	return NewLineStringZS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *LineStringZS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.LineStringZS = geo.LineStringZS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.LineStringZS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f LineStringZS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.LineStringZS)
}

// SetValid changes this LineStringZS's value and also sets it to be non-null.
func (f *LineStringZS) SetValid(n geo.LineStringZS) {
	f.LineStringZS = n
	f.Valid = true
}

// Ptr returns a pointer to this LineStringZS's value, or a nil pointer if this LineStringZS is null.
func (f LineStringZS) Ptr() *geo.LineStringZS {
	if !f.Valid {
		return nil
	}
	return &f.LineStringZS
}

// IsZero returns true for invalid LineStringZSs, for future omitempty support (Go 1.4?)
func (f LineStringZS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *LineStringZS) Scan(value interface{}) error {
	if value == nil {
		f.LineStringZS, f.Valid = geo.LineStringZS{}, false
		return nil
	}
	f.Valid = true
	return f.LineStringZS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f LineStringZS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.LineStringZS.Value()
}
