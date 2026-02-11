package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// LineStringS is a nullable geo.LineStringS.
type LineStringS struct {
	LineStringS geo.LineStringS
	Valid bool
}

// NewLineStringS creates a new LineStringS
func NewLineStringS(f geo.LineStringS, valid bool) LineStringS {
	return LineStringS{
		LineStringS: f,
		Valid: valid,
	}
}

// LineStringSFrom creates a new LineStringS that will always be valid.
func LineStringSFrom(f geo.LineStringS) LineStringS {
	return NewLineStringS(f, true)
}

// LineStringSFromPtr creates a new LineStringS that be null if f is nil.
func LineStringSFromPtr(f *geo.LineStringS) LineStringS {
	if f == nil {
		return NewLineStringS(geo.LineStringS{}, false)
	}
	return NewLineStringS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *LineStringS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.LineStringS = geo.LineStringS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.LineStringS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f LineStringS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.LineStringS)
}

// SetValid changes this LineStringS's value and also sets it to be non-null.
func (f *LineStringS) SetValid(n geo.LineStringS) {
	f.LineStringS = n
	f.Valid = true
}

// Ptr returns a pointer to this LineStringS's value, or a nil pointer if this LineStringS is null.
func (f LineStringS) Ptr() *geo.LineStringS {
	if !f.Valid {
		return nil
	}
	return &f.LineStringS
}

// IsZero returns true for invalid LineStringSs, for future omitempty support (Go 1.4?)
func (f LineStringS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *LineStringS) Scan(value interface{}) error {
	if value == nil {
		f.LineStringS, f.Valid = geo.LineStringS{}, false
		return nil
	}
	f.Valid = true
	return f.LineStringS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f LineStringS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.LineStringS.Value()
}
