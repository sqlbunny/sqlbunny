package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// MultiPointZM is a nullable geo.MultiPointZM.
type MultiPointZM struct {
	MultiPointZM geo.MultiPointZM
	Valid bool
}

// NewMultiPointZM creates a new MultiPointZM
func NewMultiPointZM(f geo.MultiPointZM, valid bool) MultiPointZM {
	return MultiPointZM{
		MultiPointZM: f,
		Valid: valid,
	}
}

// MultiPointZMFrom creates a new MultiPointZM that will always be valid.
func MultiPointZMFrom(f geo.MultiPointZM) MultiPointZM {
	return NewMultiPointZM(f, true)
}

// MultiPointZMFromPtr creates a new MultiPointZM that be null if f is nil.
func MultiPointZMFromPtr(f *geo.MultiPointZM) MultiPointZM {
	if f == nil {
		return NewMultiPointZM(geo.MultiPointZM{}, false)
	}
	return NewMultiPointZM(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPointZM) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPointZM = geo.MultiPointZM{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPointZM); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPointZM) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPointZM)
}

// SetValid changes this MultiPointZM's value and also sets it to be non-null.
func (f *MultiPointZM) SetValid(n geo.MultiPointZM) {
	f.MultiPointZM = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPointZM's value, or a nil pointer if this MultiPointZM is null.
func (f MultiPointZM) Ptr() *geo.MultiPointZM {
	if !f.Valid {
		return nil
	}
	return &f.MultiPointZM
}

// IsZero returns true for invalid MultiPointZMs, for future omitempty support (Go 1.4?)
func (f MultiPointZM) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPointZM) Scan(value interface{}) error {
	if value == nil {
		f.MultiPointZM, f.Valid = geo.MultiPointZM{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPointZM.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPointZM) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPointZM.Value()
}
