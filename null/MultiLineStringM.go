package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// MultiLineStringM is a nullable geo.MultiLineStringM.
type MultiLineStringM struct {
	MultiLineStringM geo.MultiLineStringM
	Valid bool
}

// NewMultiLineStringM creates a new MultiLineStringM
func NewMultiLineStringM(f geo.MultiLineStringM, valid bool) MultiLineStringM {
	return MultiLineStringM{
		MultiLineStringM: f,
		Valid: valid,
	}
}

// MultiLineStringMFrom creates a new MultiLineStringM that will always be valid.
func MultiLineStringMFrom(f geo.MultiLineStringM) MultiLineStringM {
	return NewMultiLineStringM(f, true)
}

// MultiLineStringMFromPtr creates a new MultiLineStringM that be null if f is nil.
func MultiLineStringMFromPtr(f *geo.MultiLineStringM) MultiLineStringM {
	if f == nil {
		return NewMultiLineStringM(geo.MultiLineStringM{}, false)
	}
	return NewMultiLineStringM(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiLineStringM) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiLineStringM = geo.MultiLineStringM{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiLineStringM); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiLineStringM) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiLineStringM)
}

// SetValid changes this MultiLineStringM's value and also sets it to be non-null.
func (f *MultiLineStringM) SetValid(n geo.MultiLineStringM) {
	f.MultiLineStringM = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiLineStringM's value, or a nil pointer if this MultiLineStringM is null.
func (f MultiLineStringM) Ptr() *geo.MultiLineStringM {
	if !f.Valid {
		return nil
	}
	return &f.MultiLineStringM
}

// IsZero returns true for invalid MultiLineStringMs, for future omitempty support (Go 1.4?)
func (f MultiLineStringM) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiLineStringM) Scan(value interface{}) error {
	if value == nil {
		f.MultiLineStringM, f.Valid = geo.MultiLineStringM{}, false
		return nil
	}
	f.Valid = true
	return f.MultiLineStringM.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiLineStringM) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiLineStringM.Value()
}
