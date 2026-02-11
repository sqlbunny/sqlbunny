package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/geo"
)

// PointZMS is a nullable geo.PointZMS.
type PointZMS struct {
	PointZMS geo.PointZMS
	Valid bool
}

// NewPointZMS creates a new PointZMS
func NewPointZMS(f geo.PointZMS, valid bool) PointZMS {
	return PointZMS{
		PointZMS: f,
		Valid: valid,
	}
}

// PointZMSFrom creates a new PointZMS that will always be valid.
func PointZMSFrom(f geo.PointZMS) PointZMS {
	return NewPointZMS(f, true)
}

// PointZMSFromPtr creates a new PointZMS that be null if f is nil.
func PointZMSFromPtr(f *geo.PointZMS) PointZMS {
	if f == nil {
		return NewPointZMS(geo.PointZMS{}, false)
	}
	return NewPointZMS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *PointZMS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.PointZMS = geo.PointZMS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.PointZMS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f PointZMS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.PointZMS)
}

// SetValid changes this PointZMS's value and also sets it to be non-null.
func (f *PointZMS) SetValid(n geo.PointZMS) {
	f.PointZMS = n
	f.Valid = true
}

// Ptr returns a pointer to this PointZMS's value, or a nil pointer if this PointZMS is null.
func (f PointZMS) Ptr() *geo.PointZMS {
	if !f.Valid {
		return nil
	}
	return &f.PointZMS
}

// IsZero returns true for invalid PointZMSs, for future omitempty support (Go 1.4?)
func (f PointZMS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *PointZMS) Scan(value interface{}) error {
	if value == nil {
		f.PointZMS, f.Valid = geo.PointZMS{}, false
		return nil
	}
	f.Valid = true
	return f.PointZMS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f PointZMS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.PointZMS.Value()
}
