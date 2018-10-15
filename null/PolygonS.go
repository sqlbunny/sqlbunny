package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// PolygonS is a nullable geo.PolygonS.
type PolygonS struct {
	PolygonS geo.PolygonS
	Valid bool
}

// NewPolygonS creates a new PolygonS
func NewPolygonS(f geo.PolygonS, valid bool) PolygonS {
	return PolygonS{
		PolygonS: f,
		Valid: valid,
	}
}

// PolygonSFrom creates a new PolygonS that will always be valid.
func PolygonSFrom(f geo.PolygonS) PolygonS {
	return NewPolygonS(f, true)
}

// PolygonSFromPtr creates a new PolygonS that be null if f is nil.
func PolygonSFromPtr(f *geo.PolygonS) PolygonS {
	if f == nil {
		return NewPolygonS(geo.PolygonS{}, false)
	}
	return NewPolygonS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *PolygonS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.PolygonS = geo.PolygonS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.PolygonS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f PolygonS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.PolygonS)
}

// SetValid changes this PolygonS's value and also sets it to be non-null.
func (f *PolygonS) SetValid(n geo.PolygonS) {
	f.PolygonS = n
	f.Valid = true
}

// Ptr returns a pointer to this PolygonS's value, or a nil pointer if this PolygonS is null.
func (f PolygonS) Ptr() *geo.PolygonS {
	if !f.Valid {
		return nil
	}
	return &f.PolygonS
}

// IsZero returns true for invalid PolygonSs, for future omitempty support (Go 1.4?)
func (f PolygonS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *PolygonS) Scan(value interface{}) error {
	if value == nil {
		f.PolygonS, f.Valid = geo.PolygonS{}, false
		return nil
	}
	f.Valid = true
	return f.PolygonS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f PolygonS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.PolygonS.Value()
}
