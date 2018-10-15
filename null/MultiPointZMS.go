package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// MultiPointZMS is a nullable geo.MultiPointZMS.
type MultiPointZMS struct {
	MultiPointZMS geo.MultiPointZMS
	Valid bool
}

// NewMultiPointZMS creates a new MultiPointZMS
func NewMultiPointZMS(f geo.MultiPointZMS, valid bool) MultiPointZMS {
	return MultiPointZMS{
		MultiPointZMS: f,
		Valid: valid,
	}
}

// MultiPointZMSFrom creates a new MultiPointZMS that will always be valid.
func MultiPointZMSFrom(f geo.MultiPointZMS) MultiPointZMS {
	return NewMultiPointZMS(f, true)
}

// MultiPointZMSFromPtr creates a new MultiPointZMS that be null if f is nil.
func MultiPointZMSFromPtr(f *geo.MultiPointZMS) MultiPointZMS {
	if f == nil {
		return NewMultiPointZMS(geo.MultiPointZMS{}, false)
	}
	return NewMultiPointZMS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPointZMS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPointZMS = geo.MultiPointZMS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPointZMS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPointZMS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPointZMS)
}

// SetValid changes this MultiPointZMS's value and also sets it to be non-null.
func (f *MultiPointZMS) SetValid(n geo.MultiPointZMS) {
	f.MultiPointZMS = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPointZMS's value, or a nil pointer if this MultiPointZMS is null.
func (f MultiPointZMS) Ptr() *geo.MultiPointZMS {
	if !f.Valid {
		return nil
	}
	return &f.MultiPointZMS
}

// IsZero returns true for invalid MultiPointZMSs, for future omitempty support (Go 1.4?)
func (f MultiPointZMS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPointZMS) Scan(value interface{}) error {
	if value == nil {
		f.MultiPointZMS, f.Valid = geo.MultiPointZMS{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPointZMS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPointZMS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPointZMS.Value()
}
