package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// MultiPoint is a nullable geo.MultiPoint.
type MultiPoint struct {
	MultiPoint geo.MultiPoint
	Valid bool
}

// NewMultiPoint creates a new MultiPoint
func NewMultiPoint(f geo.MultiPoint, valid bool) MultiPoint {
	return MultiPoint{
		MultiPoint: f,
		Valid: valid,
	}
}

// MultiPointFrom creates a new MultiPoint that will always be valid.
func MultiPointFrom(f geo.MultiPoint) MultiPoint {
	return NewMultiPoint(f, true)
}

// MultiPointFromPtr creates a new MultiPoint that be null if f is nil.
func MultiPointFromPtr(f *geo.MultiPoint) MultiPoint {
	if f == nil {
		return NewMultiPoint(geo.MultiPoint{}, false)
	}
	return NewMultiPoint(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *MultiPoint) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.MultiPoint = geo.MultiPoint{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.MultiPoint); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f MultiPoint) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.MultiPoint)
}

// SetValid changes this MultiPoint's value and also sets it to be non-null.
func (f *MultiPoint) SetValid(n geo.MultiPoint) {
	f.MultiPoint = n
	f.Valid = true
}

// Ptr returns a pointer to this MultiPoint's value, or a nil pointer if this MultiPoint is null.
func (f MultiPoint) Ptr() *geo.MultiPoint {
	if !f.Valid {
		return nil
	}
	return &f.MultiPoint
}

// IsZero returns true for invalid MultiPoints, for future omitempty support (Go 1.4?)
func (f MultiPoint) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *MultiPoint) Scan(value interface{}) error {
	if value == nil {
		f.MultiPoint, f.Valid = geo.MultiPoint{}, false
		return nil
	}
	f.Valid = true
	return f.MultiPoint.Scan(value)
}

// Value implements the driver Valuer interface.
func (f MultiPoint) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.MultiPoint.Value()
}
