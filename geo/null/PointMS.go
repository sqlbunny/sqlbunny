package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/geo"
)

// PointMS is a nullable geo.PointMS.
type PointMS struct {
	PointMS geo.PointMS
	Valid bool
}

// NewPointMS creates a new PointMS
func NewPointMS(f geo.PointMS, valid bool) PointMS {
	return PointMS{
		PointMS: f,
		Valid: valid,
	}
}

// PointMSFrom creates a new PointMS that will always be valid.
func PointMSFrom(f geo.PointMS) PointMS {
	return NewPointMS(f, true)
}

// PointMSFromPtr creates a new PointMS that be null if f is nil.
func PointMSFromPtr(f *geo.PointMS) PointMS {
	if f == nil {
		return NewPointMS(geo.PointMS{}, false)
	}
	return NewPointMS(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *PointMS) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.PointMS = geo.PointMS{}
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.PointMS); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f PointMS) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return json.Marshal(f.PointMS)
}

// SetValid changes this PointMS's value and also sets it to be non-null.
func (f *PointMS) SetValid(n geo.PointMS) {
	f.PointMS = n
	f.Valid = true
}

// Ptr returns a pointer to this PointMS's value, or a nil pointer if this PointMS is null.
func (f PointMS) Ptr() *geo.PointMS {
	if !f.Valid {
		return nil
	}
	return &f.PointMS
}

// IsZero returns true for invalid PointMSs, for future omitempty support (Go 1.4?)
func (f PointMS) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *PointMS) Scan(value interface{}) error {
	if value == nil {
		f.PointMS, f.Valid = geo.PointMS{}, false
		return nil
	}
	f.Valid = true
	return f.PointMS.Scan(value)
}

// Value implements the driver Valuer interface.
func (f PointMS) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.PointMS.Value()
}
