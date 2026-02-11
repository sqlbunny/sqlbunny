package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/geo"
)

// Point is a nullable geo.Point.
type Point struct {
	Point geo.Point
	Valid bool
}

// NewPoint creates a new Point
func NewPoint(f geo.Point, valid bool) Point {
	return Point{
		Point: f,
		Valid: valid,
	}
}

// PointFrom creates a new Point that will always be valid.
func PointFrom(f geo.Point) Point {
	return NewPoint(f, true)
}

// PointFromPtr creates a new Point that be null if f is nil.
func PointFromPtr(f *geo.Point) Point {
	if f == nil {
		return NewPoint(geo.Point{}, false)
	}
	return NewPoint(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Point) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.Point = geo.Point{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.Point); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f Point) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.Point)
}

// SetValid changes this Point's value and also sets it to be non-null.
func (f *Point) SetValid(n geo.Point) {
	f.Point = n
	f.Valid = true
}

// Ptr returns a pointer to this Point's value, or a nil pointer if this Point is null.
func (f Point) Ptr() *geo.Point {
	if !f.Valid {
		return nil
	}
	return &f.Point
}

// IsZero returns true for invalid Points, for future omitempty support (Go 1.4?)
func (f Point) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *Point) Scan(value interface{}) error {
	if value == nil {
		f.Point, f.Valid = geo.Point{}, false
		return nil
	}
	f.Valid = true
	return f.Point.Scan(value)
}

// Value implements the driver Valuer interface.
func (f Point) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.Point.Value()
}
