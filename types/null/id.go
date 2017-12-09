package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/volatiletech/sqlboiler/types"
)

// ID is an nullable ID.
type ID struct {
	ID    types.ID
	Valid bool
}

// NewID creates a new ID
func NewID(i types.ID, valid bool) ID {
	return ID{
		ID:    i,
		Valid: valid,
	}
}

// IDFrom creates a new ID that will always be valid.
func IDFrom(i types.ID) ID {
	return NewID(i, true)
}

// IDFromPtr creates a new ID that be null if i is nil.
func IDFromPtr(i *types.ID) ID {
	if i == nil {
		return NewID(types.ID{}, false)
	}
	return NewID(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *ID) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		u.ID = types.ID{}
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.ID); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *ID) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := types.IDFromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.ID = res
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (u ID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return NullBytes, nil
	}
	return u.ID.MarshalJSON()
}

// MarshalText implements encoding.TextMarshaler.
func (u ID) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return u.ID.MarshalText()
}

// SetValid changes this ID's value and also sets it to be non-null.
func (u *ID) SetValid(n types.ID) {
	u.ID = n
	u.Valid = true
}

// Ptr returns a pointer to this ID's value, or a nil pointer if this ID is null.
func (u ID) Ptr() *types.ID {
	if !u.Valid {
		return nil
	}
	return &u.ID
}

// IsZero returns true for invalid ID's, for future omitempty support (Go 1.4?)
func (u ID) IsZero() bool {
	return !u.Valid
}

// Scan implements the Scanner interface.
func (u *ID) Scan(value interface{}) error {
	if value == nil {
		u.ID, u.Valid = types.ID{}, false
		return nil
	}
	u.Valid = true
	return u.ID.Scan(value)
}

// Value implements the driver Valuer interface.
func (u ID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.Value()
}
