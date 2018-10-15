package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/kernelpayments/geo"
)

// MultiPolygonS is a nullable geo.MultiPolygonS.
type MultiPolygonS struct {
	MultiPolygonS geo.MultiPolygonS
	Valid bool
}

// NewMultiPolygonS creates a new MultiPolygonS
func NewMultiPolygonS(f geo.MultiPolygonS, valid bool) MultiPolygonS {
	return MultiPolygonS{
		MultiPolygonS: f,
		Valid: valid,
	}
}

// MultiPolygonSFrom creates a new MultiPolygonS that will always be valid.
func MultiPolygonSFrom(f geo.MultiPolygonS) MultiPolygonS {
	return NewMultiPolygonS(f, true)
}

// MultiPolygonSFromPtr creates a new MultiPolygonS that be null if f is nil.
func MultiPolygonSFromPtr(f *geo.MultiPolygonS) MultiPolygonS {
	if f == nil {
		return NewMultiPolygonS(geo.MultiPolygonS{}, false)
	}
	return NewMultiPolygonS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPolygonS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPolygonS = geo.MultiPolygonS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPolygonS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPolygonS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPolygonS)
}

// SetValid changes this MultiPolygonS's value and also sets it to be non-null.
func (f *MultiPolygonS) SetValid(n geo.MultiPolygonS) {
	f.MultiPolygonS = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPolygonS's value, or a nil pointer if this MultiPolygonS is null.
func (f MultiPolygonS) Ptr() *geo.MultiPolygonS {
	if !f.Valid {
		return nil
	}
	return &f.MultiPolygonS
}

// IsZero returns true for invalid MultiPolygonSs, for future omitempty support (Go 1.4?)
func (f MultiPolygonS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPolygonS) Scan(value interface{}) error {
	if value == nil {
		f.MultiPolygonS, f.Valid = geo.MultiPolygonS{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPolygonS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPolygonS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPolygonS.Value()
}
