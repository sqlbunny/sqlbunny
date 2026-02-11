package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/geo"
)

// PolygonZS is a nullable geo.PolygonZS.
type PolygonZS struct {
	PolygonZS geo.PolygonZS
	Valid bool
}

// NewPolygonZS creates a new PolygonZS
func NewPolygonZS(f geo.PolygonZS, valid bool) PolygonZS {
	return PolygonZS{
		PolygonZS: f,
		Valid: valid,
	}
}

// PolygonZSFrom creates a new PolygonZS that will always be valid.
func PolygonZSFrom(f geo.PolygonZS) PolygonZS {
	return NewPolygonZS(f, true)
}

// PolygonZSFromPtr creates a new PolygonZS that be null if f is nil.
func PolygonZSFromPtr(f *geo.PolygonZS) PolygonZS {
	if f == nil {
		return NewPolygonZS(geo.PolygonZS{}, false)
	}
	return NewPolygonZS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *PolygonZS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.PolygonZS = geo.PolygonZS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.PolygonZS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f PolygonZS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.PolygonZS)
}

// SetValid changes this PolygonZS's value and also sets it to be non-null.
func (f *PolygonZS) SetValid(n geo.PolygonZS) {
	f.PolygonZS = n
	f.Valid = true
}

// Ptr returns a pointer to this PolygonZS's value, or a nil pointer if this PolygonZS is null.
func (f PolygonZS) Ptr() *geo.PolygonZS {
	if !f.Valid {
		return nil
	}
	return &f.PolygonZS
}

// IsZero returns true for invalid PolygonZSs, for future omitempty support (Go 1.4?)
func (f PolygonZS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *PolygonZS) Scan(value interface{}) error {
	if value == nil {
		f.PolygonZS, f.Valid = geo.PolygonZS{}, false
		return nil
	}
	f.Valid = true
	return f.PolygonZS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f PolygonZS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.PolygonZS.Value()
}
