package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// MultiLineStringZM is a nullable geo.MultiLineStringZM.
type MultiLineStringZM struct {
	MultiLineStringZM geo.MultiLineStringZM
	Valid bool
}

// NewMultiLineStringZM creates a new MultiLineStringZM
func NewMultiLineStringZM(f geo.MultiLineStringZM, valid bool) MultiLineStringZM {
	return MultiLineStringZM{
		MultiLineStringZM: f,
		Valid: valid,
	}
}

// MultiLineStringZMFrom creates a new MultiLineStringZM that will always be valid.
func MultiLineStringZMFrom(f geo.MultiLineStringZM) MultiLineStringZM {
	return NewMultiLineStringZM(f, true)
}

// MultiLineStringZMFromPtr creates a new MultiLineStringZM that be null if f is nil.
func MultiLineStringZMFromPtr(f *geo.MultiLineStringZM) MultiLineStringZM {
	if f == nil {
		return NewMultiLineStringZM(geo.MultiLineStringZM{}, false)
	}
	return NewMultiLineStringZM(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiLineStringZM) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiLineStringZM = geo.MultiLineStringZM{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiLineStringZM); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiLineStringZM) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiLineStringZM)
}

// SetValid changes this MultiLineStringZM's value and also sets it to be non-null.
func (f *MultiLineStringZM) SetValid(n geo.MultiLineStringZM) {
	f.MultiLineStringZM = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiLineStringZM's value, or a nil pointer if this MultiLineStringZM is null.
func (f MultiLineStringZM) Ptr() *geo.MultiLineStringZM {
	if !f.Valid {
		return nil
	}
	return &f.MultiLineStringZM
}

// IsZero returns true for invalid MultiLineStringZMs, for future omitempty support (Go 1.4?)
func (f MultiLineStringZM) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiLineStringZM) Scan(value interface{}) error {
	if value == nil {
		f.MultiLineStringZM, f.Valid = geo.MultiLineStringZM{}, false
		return nil
	}
	f.Valid = true
	return f.MultiLineStringZM.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiLineStringZM) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiLineStringZM.Value()
}
