package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// PolygonM is a nullable geo.PolygonM.
type PolygonM struct {
	PolygonM geo.PolygonM
	Valid bool
}

// NewPolygonM creates a new PolygonM
func NewPolygonM(f geo.PolygonM, valid bool) PolygonM {
	return PolygonM{
		PolygonM: f,
		Valid: valid,
	}
}

// PolygonMFrom creates a new PolygonM that will always be valid.
func PolygonMFrom(f geo.PolygonM) PolygonM {
	return NewPolygonM(f, true)
}

// PolygonMFromPtr creates a new PolygonM that be null if f is nil.
func PolygonMFromPtr(f *geo.PolygonM) PolygonM {
	if f == nil {
		return NewPolygonM(geo.PolygonM{}, false)
	}
	return NewPolygonM(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *PolygonM) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.PolygonM = geo.PolygonM{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.PolygonM); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f PolygonM) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.PolygonM)
}

// SetValid changes this PolygonM's value and also sets it to be non-null.
func (f *PolygonM) SetValid(n geo.PolygonM) {
	f.PolygonM = n
	f.Valid = true
}

// Ptr returns a pointer to this PolygonM's value, or a nil pointer if this PolygonM is null.
func (f PolygonM) Ptr() *geo.PolygonM {
	if !f.Valid {
		return nil
	}
	return &f.PolygonM
}

// IsZero returns true for invalid PolygonMs, for future omitempty support (Go 1.4?)
func (f PolygonM) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *PolygonM) Scan(value interface{}) error {
	if value == nil {
		f.PolygonM, f.Valid = geo.PolygonM{}, false
		return nil
	}
	f.Valid = true
	return f.PolygonM.Scan(value)
}

// Value implements the driver Valuer interface.
func (f PolygonM) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.PolygonM.Value()
}
