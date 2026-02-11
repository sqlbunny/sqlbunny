package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// MultiPointZS is a nullable geo.MultiPointZS.
type MultiPointZS struct {
	MultiPointZS geo.MultiPointZS
	Valid bool
}

// NewMultiPointZS creates a new MultiPointZS
func NewMultiPointZS(f geo.MultiPointZS, valid bool) MultiPointZS {
	return MultiPointZS{
		MultiPointZS: f,
		Valid: valid,
	}
}

// MultiPointZSFrom creates a new MultiPointZS that will always be valid.
func MultiPointZSFrom(f geo.MultiPointZS) MultiPointZS {
	return NewMultiPointZS(f, true)
}

// MultiPointZSFromPtr creates a new MultiPointZS that be null if f is nil.
func MultiPointZSFromPtr(f *geo.MultiPointZS) MultiPointZS {
	if f == nil {
		return NewMultiPointZS(geo.MultiPointZS{}, false)
	}
	return NewMultiPointZS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPointZS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPointZS = geo.MultiPointZS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPointZS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPointZS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPointZS)
}

// SetValid changes this MultiPointZS's value and also sets it to be non-null.
func (f *MultiPointZS) SetValid(n geo.MultiPointZS) {
	f.MultiPointZS = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPointZS's value, or a nil pointer if this MultiPointZS is null.
func (f MultiPointZS) Ptr() *geo.MultiPointZS {
	if !f.Valid {
		return nil
	}
	return &f.MultiPointZS
}

// IsZero returns true for invalid MultiPointZSs, for future omitempty support (Go 1.4?)
func (f MultiPointZS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPointZS) Scan(value interface{}) error {
	if value == nil {
		f.MultiPointZS, f.Valid = geo.MultiPointZS{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPointZS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPointZS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPointZS.Value()
}
