package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// MultiPolygonZS is a nullable geo.MultiPolygonZS.
type MultiPolygonZS struct {
	MultiPolygonZS geo.MultiPolygonZS
	Valid bool
}

// NewMultiPolygonZS creates a new MultiPolygonZS
func NewMultiPolygonZS(f geo.MultiPolygonZS, valid bool) MultiPolygonZS {
	return MultiPolygonZS{
		MultiPolygonZS: f,
		Valid: valid,
	}
}

// MultiPolygonZSFrom creates a new MultiPolygonZS that will always be valid.
func MultiPolygonZSFrom(f geo.MultiPolygonZS) MultiPolygonZS {
	return NewMultiPolygonZS(f, true)
}

// MultiPolygonZSFromPtr creates a new MultiPolygonZS that be null if f is nil.
func MultiPolygonZSFromPtr(f *geo.MultiPolygonZS) MultiPolygonZS {
	if f == nil {
		return NewMultiPolygonZS(geo.MultiPolygonZS{}, false)
	}
	return NewMultiPolygonZS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPolygonZS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPolygonZS = geo.MultiPolygonZS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPolygonZS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPolygonZS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPolygonZS)
}

// SetValid changes this MultiPolygonZS's value and also sets it to be non-null.
func (f *MultiPolygonZS) SetValid(n geo.MultiPolygonZS) {
	f.MultiPolygonZS = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPolygonZS's value, or a nil pointer if this MultiPolygonZS is null.
func (f MultiPolygonZS) Ptr() *geo.MultiPolygonZS {
	if !f.Valid {
		return nil
	}
	return &f.MultiPolygonZS
}

// IsZero returns true for invalid MultiPolygonZSs, for future omitempty support (Go 1.4?)
func (f MultiPolygonZS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPolygonZS) Scan(value interface{}) error {
	if value == nil {
		f.MultiPolygonZS, f.Valid = geo.MultiPolygonZS{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPolygonZS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPolygonZS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPolygonZS.Value()
}
