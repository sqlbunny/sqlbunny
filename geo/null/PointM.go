package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/geo"
)

// PointM is a nullable geo.PointM.
type PointM struct {
	PointM geo.PointM
	Valid bool
}

// NewPointM creates a new PointM
func NewPointM(f geo.PointM, valid bool) PointM {
	return PointM{
		PointM: f,
		Valid: valid,
	}
}

// PointMFrom creates a new PointM that will always be valid.
func PointMFrom(f geo.PointM) PointM {
	return NewPointM(f, true)
}

// PointMFromPtr creates a new PointM that be null if f is nil.
func PointMFromPtr(f *geo.PointM) PointM {
	if f == nil {
		return NewPointM(geo.PointM{}, false)
	}
	return NewPointM(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *PointM) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.PointM = geo.PointM{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.PointM); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f PointM) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.PointM)
}

// SetValid changes this PointM's value and also sets it to be non-null.
func (f *PointM) SetValid(n geo.PointM) {
	f.PointM = n
	f.Valid = true
}

// Ptr returns a pointer to this PointM's value, or a nil pointer if this PointM is null.
func (f PointM) Ptr() *geo.PointM {
	if !f.Valid {
		return nil
	}
	return &f.PointM
}

// IsZero returns true for invalid PointMs, for future omitempty support (Go 1.4?)
func (f PointM) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *PointM) Scan(value interface{}) error {
	if value == nil {
		f.PointM, f.Valid = geo.PointM{}, false
		return nil
	}
	f.Valid = true
	return f.PointM.Scan(value)
}

// Value implements the driver Valuer interface.
func (f PointM) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.PointM.Value()
}
