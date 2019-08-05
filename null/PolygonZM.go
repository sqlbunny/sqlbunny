package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// PolygonZM is a nullable geo.PolygonZM.
type PolygonZM struct {
	PolygonZM geo.PolygonZM
	Valid bool
}

// NewPolygonZM creates a new PolygonZM
func NewPolygonZM(f geo.PolygonZM, valid bool) PolygonZM {
	return PolygonZM{
		PolygonZM: f,
		Valid: valid,
	}
}

// PolygonZMFrom creates a new PolygonZM that will always be valid.
func PolygonZMFrom(f geo.PolygonZM) PolygonZM {
	return NewPolygonZM(f, true)
}

// PolygonZMFromPtr creates a new PolygonZM that be null if f is nil.
func PolygonZMFromPtr(f *geo.PolygonZM) PolygonZM {
	if f == nil {
		return NewPolygonZM(geo.PolygonZM{}, false)
	}
	return NewPolygonZM(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *PolygonZM) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.PolygonZM = geo.PolygonZM{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.PolygonZM); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f PolygonZM) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.PolygonZM)
}

// SetValid changes this PolygonZM's value and also sets it to be non-null.
func (f *PolygonZM) SetValid(n geo.PolygonZM) {
	f.PolygonZM = n
	f.Valid = true
}

// Ptr returns a pointer to this PolygonZM's value, or a nil pointer if this PolygonZM is null.
func (f PolygonZM) Ptr() *geo.PolygonZM {
	if !f.Valid {
		return nil
	}
	return &f.PolygonZM
}

// IsZero returns true for invalid PolygonZMs, for future omitempty support (Go 1.4?)
func (f PolygonZM) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *PolygonZM) Scan(value interface{}) error {
	if value == nil {
		f.PolygonZM, f.Valid = geo.PolygonZM{}, false
		return nil
	}
	f.Valid = true
	return f.PolygonZM.Scan(value)
}

// Value implements the driver Valuer interface.
func (f PolygonZM) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.PolygonZM.Value()
}
