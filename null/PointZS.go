package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// PointZS is a nullable geo.PointZS.
type PointZS struct {
	PointZS geo.PointZS
	Valid bool
}

// NewPointZS creates a new PointZS
func NewPointZS(f geo.PointZS, valid bool) PointZS {
	return PointZS{
		PointZS: f,
		Valid: valid,
	}
}

// PointZSFrom creates a new PointZS that will always be valid.
func PointZSFrom(f geo.PointZS) PointZS {
	return NewPointZS(f, true)
}

// PointZSFromPtr creates a new PointZS that be null if f is nil.
func PointZSFromPtr(f *geo.PointZS) PointZS {
	if f == nil {
		return NewPointZS(geo.PointZS{}, false)
	}
	return NewPointZS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *PointZS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.PointZS = geo.PointZS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.PointZS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f PointZS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.PointZS)
}

// SetValid changes this PointZS's value and also sets it to be non-null.
func (f *PointZS) SetValid(n geo.PointZS) {
	f.PointZS = n
	f.Valid = true
}

// Ptr returns a pointer to this PointZS's value, or a nil pointer if this PointZS is null.
func (f PointZS) Ptr() *geo.PointZS {
	if !f.Valid {
		return nil
	}
	return &f.PointZS
}

// IsZero returns true for invalid PointZSs, for future omitempty support (Go 1.4?)
func (f PointZS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *PointZS) Scan(value interface{}) error {
	if value == nil {
		f.PointZS, f.Valid = geo.PointZS{}, false
		return nil
	}
	f.Valid = true
	return f.PointZS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f PointZS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.PointZS.Value()
}
