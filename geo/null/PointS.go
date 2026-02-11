package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/geo"
)

// PointS is a nullable geo.PointS.
type PointS struct {
	PointS geo.PointS
	Valid bool
}

// NewPointS creates a new PointS
func NewPointS(f geo.PointS, valid bool) PointS {
	return PointS{
		PointS: f,
		Valid: valid,
	}
}

// PointSFrom creates a new PointS that will always be valid.
func PointSFrom(f geo.PointS) PointS {
	return NewPointS(f, true)
}

// PointSFromPtr creates a new PointS that be null if f is nil.
func PointSFromPtr(f *geo.PointS) PointS {
	if f == nil {
		return NewPointS(geo.PointS{}, false)
	}
	return NewPointS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *PointS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.PointS = geo.PointS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.PointS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f PointS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.PointS)
}

// SetValid changes this PointS's value and also sets it to be non-null.
func (f *PointS) SetValid(n geo.PointS) {
	f.PointS = n
	f.Valid = true
}

// Ptr returns a pointer to this PointS's value, or a nil pointer if this PointS is null.
func (f PointS) Ptr() *geo.PointS {
	if !f.Valid {
		return nil
	}
	return &f.PointS
}

// IsZero returns true for invalid PointSs, for future omitempty support (Go 1.4?)
func (f PointS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *PointS) Scan(value interface{}) error {
	if value == nil {
		f.PointS, f.Valid = geo.PointS{}, false
		return nil
	}
	f.Valid = true
	return f.PointS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f PointS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.PointS.Value()
}
