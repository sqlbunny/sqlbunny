package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/geo"
)

// Polygon is a nullable geo.Polygon.
type Polygon struct {
	Polygon geo.Polygon
	Valid bool
}

// NewPolygon creates a new Polygon
func NewPolygon(f geo.Polygon, valid bool) Polygon {
	return Polygon{
		Polygon: f,
		Valid: valid,
	}
}

// PolygonFrom creates a new Polygon that will always be valid.
func PolygonFrom(f geo.Polygon) Polygon {
	return NewPolygon(f, true)
}

// PolygonFromPtr creates a new Polygon that be null if f is nil.
func PolygonFromPtr(f *geo.Polygon) Polygon {
	if f == nil {
		return NewPolygon(geo.Polygon{}, false)
	}
	return NewPolygon(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Polygon) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.Polygon = geo.Polygon{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.Polygon); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f Polygon) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.Polygon)
}

// SetValid changes this Polygon's value and also sets it to be non-null.
func (f *Polygon) SetValid(n geo.Polygon) {
	f.Polygon = n
	f.Valid = true
}

// Ptr returns a pointer to this Polygon's value, or a nil pointer if this Polygon is null.
func (f Polygon) Ptr() *geo.Polygon {
	if !f.Valid {
		return nil
	}
	return &f.Polygon
}

// IsZero returns true for invalid Polygons, for future omitempty support (Go 1.4?)
func (f Polygon) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *Polygon) Scan(value interface{}) error {
	if value == nil {
		f.Polygon, f.Valid = geo.Polygon{}, false
		return nil
	}
	f.Valid = true
	return f.Polygon.Scan(value)
}

// Value implements the driver Valuer interface.
func (f Polygon) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.Polygon.Value()
}
