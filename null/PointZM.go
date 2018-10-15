package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// PointZM is a nullable geo.PointZM.
type PointZM struct {
	PointZM geo.PointZM
	Valid bool
}

// NewPointZM creates a new PointZM
func NewPointZM(f geo.PointZM, valid bool) PointZM {
	return PointZM{
		PointZM: f,
		Valid: valid,
	}
}

// PointZMFrom creates a new PointZM that will always be valid.
func PointZMFrom(f geo.PointZM) PointZM {
	return NewPointZM(f, true)
}

// PointZMFromPtr creates a new PointZM that be null if f is nil.
func PointZMFromPtr(f *geo.PointZM) PointZM {
	if f == nil {
		return NewPointZM(geo.PointZM{}, false)
	}
	return NewPointZM(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *PointZM) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.PointZM = geo.PointZM{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.PointZM); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f PointZM) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.PointZM)
}

// SetValid changes this PointZM's value and also sets it to be non-null.
func (f *PointZM) SetValid(n geo.PointZM) {
	f.PointZM = n
	f.Valid = true
}

// Ptr returns a pointer to this PointZM's value, or a nil pointer if this PointZM is null.
func (f PointZM) Ptr() *geo.PointZM {
	if !f.Valid {
		return nil
	}
	return &f.PointZM
}

// IsZero returns true for invalid PointZMs, for future omitempty support (Go 1.4?)
func (f PointZM) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *PointZM) Scan(value interface{}) error {
	if value == nil {
		f.PointZM, f.Valid = geo.PointZM{}, false
		return nil
	}
	f.Valid = true
	return f.PointZM.Scan(value)
}

// Value implements the driver Valuer interface.
func (f PointZM) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.PointZM.Value()
}
