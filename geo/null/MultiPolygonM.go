package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/geo"
)

// MultiPolygonM is a nullable geo.MultiPolygonM.
type MultiPolygonM struct {
	MultiPolygonM geo.MultiPolygonM
	Valid bool
}

// NewMultiPolygonM creates a new MultiPolygonM
func NewMultiPolygonM(f geo.MultiPolygonM, valid bool) MultiPolygonM {
	return MultiPolygonM{
		MultiPolygonM: f,
		Valid: valid,
	}
}

// MultiPolygonMFrom creates a new MultiPolygonM that will always be valid.
func MultiPolygonMFrom(f geo.MultiPolygonM) MultiPolygonM {
	return NewMultiPolygonM(f, true)
}

// MultiPolygonMFromPtr creates a new MultiPolygonM that be null if f is nil.
func MultiPolygonMFromPtr(f *geo.MultiPolygonM) MultiPolygonM {
	if f == nil {
		return NewMultiPolygonM(geo.MultiPolygonM{}, false)
	}
	return NewMultiPolygonM(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPolygonM) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPolygonM = geo.MultiPolygonM{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPolygonM); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPolygonM) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPolygonM)
}

// SetValid changes this MultiPolygonM's value and also sets it to be non-null.
func (f *MultiPolygonM) SetValid(n geo.MultiPolygonM) {
	f.MultiPolygonM = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPolygonM's value, or a nil pointer if this MultiPolygonM is null.
func (f MultiPolygonM) Ptr() *geo.MultiPolygonM {
	if !f.Valid {
		return nil
	}
	return &f.MultiPolygonM
}

// IsZero returns true for invalid MultiPolygonMs, for future omitempty support (Go 1.4?)
func (f MultiPolygonM) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPolygonM) Scan(value interface{}) error {
	if value == nil {
		f.MultiPolygonM, f.Valid = geo.MultiPolygonM{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPolygonM.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPolygonM) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPolygonM.Value()
}
