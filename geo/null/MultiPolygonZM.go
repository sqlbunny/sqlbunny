package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/geo"
)

// MultiPolygonZM is a nullable geo.MultiPolygonZM.
type MultiPolygonZM struct {
	MultiPolygonZM geo.MultiPolygonZM
	Valid bool
}

// NewMultiPolygonZM creates a new MultiPolygonZM
func NewMultiPolygonZM(f geo.MultiPolygonZM, valid bool) MultiPolygonZM {
	return MultiPolygonZM{
		MultiPolygonZM: f,
		Valid: valid,
	}
}

// MultiPolygonZMFrom creates a new MultiPolygonZM that will always be valid.
func MultiPolygonZMFrom(f geo.MultiPolygonZM) MultiPolygonZM {
	return NewMultiPolygonZM(f, true)
}

// MultiPolygonZMFromPtr creates a new MultiPolygonZM that be null if f is nil.
func MultiPolygonZMFromPtr(f *geo.MultiPolygonZM) MultiPolygonZM {
	if f == nil {
		return NewMultiPolygonZM(geo.MultiPolygonZM{}, false)
	}
	return NewMultiPolygonZM(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPolygonZM) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPolygonZM = geo.MultiPolygonZM{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPolygonZM); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPolygonZM) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPolygonZM)
}

// SetValid changes this MultiPolygonZM's value and also sets it to be non-null.
func (f *MultiPolygonZM) SetValid(n geo.MultiPolygonZM) {
	f.MultiPolygonZM = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPolygonZM's value, or a nil pointer if this MultiPolygonZM is null.
func (f MultiPolygonZM) Ptr() *geo.MultiPolygonZM {
	if !f.Valid {
		return nil
	}
	return &f.MultiPolygonZM
}

// IsZero returns true for invalid MultiPolygonZMs, for future omitempty support (Go 1.4?)
func (f MultiPolygonZM) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPolygonZM) Scan(value interface{}) error {
	if value == nil {
		f.MultiPolygonZM, f.Valid = geo.MultiPolygonZM{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPolygonZM.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPolygonZM) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPolygonZM.Value()
}
