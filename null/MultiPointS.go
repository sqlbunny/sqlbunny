package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// MultiPointS is a nullable geo.MultiPointS.
type MultiPointS struct {
	MultiPointS geo.MultiPointS
	Valid bool
}

// NewMultiPointS creates a new MultiPointS
func NewMultiPointS(f geo.MultiPointS, valid bool) MultiPointS {
	return MultiPointS{
		MultiPointS: f,
		Valid: valid,
	}
}

// MultiPointSFrom creates a new MultiPointS that will always be valid.
func MultiPointSFrom(f geo.MultiPointS) MultiPointS {
	return NewMultiPointS(f, true)
}

// MultiPointSFromPtr creates a new MultiPointS that be null if f is nil.
func MultiPointSFromPtr(f *geo.MultiPointS) MultiPointS {
	if f == nil {
		return NewMultiPointS(geo.MultiPointS{}, false)
	}
	return NewMultiPointS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPointS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPointS = geo.MultiPointS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPointS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPointS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPointS)
}

// SetValid changes this MultiPointS's value and also sets it to be non-null.
func (f *MultiPointS) SetValid(n geo.MultiPointS) {
	f.MultiPointS = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPointS's value, or a nil pointer if this MultiPointS is null.
func (f MultiPointS) Ptr() *geo.MultiPointS {
	if !f.Valid {
		return nil
	}
	return &f.MultiPointS
}

// IsZero returns true for invalid MultiPointSs, for future omitempty support (Go 1.4?)
func (f MultiPointS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPointS) Scan(value interface{}) error {
	if value == nil {
		f.MultiPointS, f.Valid = geo.MultiPointS{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPointS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPointS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPointS.Value()
}
