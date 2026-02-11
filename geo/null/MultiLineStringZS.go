package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// MultiLineStringZS is a nullable geo.MultiLineStringZS.
type MultiLineStringZS struct {
	MultiLineStringZS geo.MultiLineStringZS
	Valid bool
}

// NewMultiLineStringZS creates a new MultiLineStringZS
func NewMultiLineStringZS(f geo.MultiLineStringZS, valid bool) MultiLineStringZS {
	return MultiLineStringZS{
		MultiLineStringZS: f,
		Valid: valid,
	}
}

// MultiLineStringZSFrom creates a new MultiLineStringZS that will always be valid.
func MultiLineStringZSFrom(f geo.MultiLineStringZS) MultiLineStringZS {
	return NewMultiLineStringZS(f, true)
}

// MultiLineStringZSFromPtr creates a new MultiLineStringZS that be null if f is nil.
func MultiLineStringZSFromPtr(f *geo.MultiLineStringZS) MultiLineStringZS {
	if f == nil {
		return NewMultiLineStringZS(geo.MultiLineStringZS{}, false)
	}
	return NewMultiLineStringZS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiLineStringZS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiLineStringZS = geo.MultiLineStringZS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiLineStringZS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiLineStringZS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiLineStringZS)
}

// SetValid changes this MultiLineStringZS's value and also sets it to be non-null.
func (f *MultiLineStringZS) SetValid(n geo.MultiLineStringZS) {
	f.MultiLineStringZS = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiLineStringZS's value, or a nil pointer if this MultiLineStringZS is null.
func (f MultiLineStringZS) Ptr() *geo.MultiLineStringZS {
	if !f.Valid {
		return nil
	}
	return &f.MultiLineStringZS
}

// IsZero returns true for invalid MultiLineStringZSs, for future omitempty support (Go 1.4?)
func (f MultiLineStringZS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiLineStringZS) Scan(value interface{}) error {
	if value == nil {
		f.MultiLineStringZS, f.Valid = geo.MultiLineStringZS{}, false
		return nil
	}
	f.Valid = true
	return f.MultiLineStringZS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiLineStringZS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiLineStringZS.Value()
}
