package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// MultiLineStringS is a nullable geo.MultiLineStringS.
type MultiLineStringS struct {
	MultiLineStringS geo.MultiLineStringS
	Valid bool
}

// NewMultiLineStringS creates a new MultiLineStringS
func NewMultiLineStringS(f geo.MultiLineStringS, valid bool) MultiLineStringS {
	return MultiLineStringS{
		MultiLineStringS: f,
		Valid: valid,
	}
}

// MultiLineStringSFrom creates a new MultiLineStringS that will always be valid.
func MultiLineStringSFrom(f geo.MultiLineStringS) MultiLineStringS {
	return NewMultiLineStringS(f, true)
}

// MultiLineStringSFromPtr creates a new MultiLineStringS that be null if f is nil.
func MultiLineStringSFromPtr(f *geo.MultiLineStringS) MultiLineStringS {
	if f == nil {
		return NewMultiLineStringS(geo.MultiLineStringS{}, false)
	}
	return NewMultiLineStringS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiLineStringS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiLineStringS = geo.MultiLineStringS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiLineStringS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiLineStringS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiLineStringS)
}

// SetValid changes this MultiLineStringS's value and also sets it to be non-null.
func (f *MultiLineStringS) SetValid(n geo.MultiLineStringS) {
	f.MultiLineStringS = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiLineStringS's value, or a nil pointer if this MultiLineStringS is null.
func (f MultiLineStringS) Ptr() *geo.MultiLineStringS {
	if !f.Valid {
		return nil
	}
	return &f.MultiLineStringS
}

// IsZero returns true for invalid MultiLineStringSs, for future omitempty support (Go 1.4?)
func (f MultiLineStringS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiLineStringS) Scan(value interface{}) error {
	if value == nil {
		f.MultiLineStringS, f.Valid = geo.MultiLineStringS{}, false
		return nil
	}
	f.Valid = true
	return f.MultiLineStringS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiLineStringS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiLineStringS.Value()
}
