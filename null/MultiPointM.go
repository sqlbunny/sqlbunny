package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// MultiPointM is a nullable geo.MultiPointM.
type MultiPointM struct {
	MultiPointM geo.MultiPointM
	Valid bool
}

// NewMultiPointM creates a new MultiPointM
func NewMultiPointM(f geo.MultiPointM, valid bool) MultiPointM {
	return MultiPointM{
		MultiPointM: f,
		Valid: valid,
	}
}

// MultiPointMFrom creates a new MultiPointM that will always be valid.
func MultiPointMFrom(f geo.MultiPointM) MultiPointM {
	return NewMultiPointM(f, true)
}

// MultiPointMFromPtr creates a new MultiPointM that be null if f is nil.
func MultiPointMFromPtr(f *geo.MultiPointM) MultiPointM {
	if f == nil {
		return NewMultiPointM(geo.MultiPointM{}, false)
	}
	return NewMultiPointM(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPointM) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPointM = geo.MultiPointM{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPointM); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPointM) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPointM)
}

// SetValid changes this MultiPointM's value and also sets it to be non-null.
func (f *MultiPointM) SetValid(n geo.MultiPointM) {
	f.MultiPointM = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPointM's value, or a nil pointer if this MultiPointM is null.
func (f MultiPointM) Ptr() *geo.MultiPointM {
	if !f.Valid {
		return nil
	}
	return &f.MultiPointM
}

// IsZero returns true for invalid MultiPointMs, for future omitempty support (Go 1.4?)
func (f MultiPointM) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPointM) Scan(value interface{}) error {
	if value == nil {
		f.MultiPointM, f.Valid = geo.MultiPointM{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPointM.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPointM) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPointM.Value()
}
