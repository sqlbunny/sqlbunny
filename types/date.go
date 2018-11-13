package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

const RFC3339Date = "2006-01-02"

// Date is a time.Time designed to store only the date. It supports SQL and JSON serialization.
type Date time.Time

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in RFC 3339 format, with sub-second precision added if present.
func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Date.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(RFC3339Date)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, RFC3339Date)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (d *Date) UnmarshalJSON(data []byte) error {
	// Fractional seconds are handled implicitly by Parse.
	t, err := time.Parse(`"`+RFC3339Date+`"`, string(data))
	if err != nil {
		return err
	}
	*d = Date(t)
	return err
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted in RFC 3339 format, with sub-second precision added if present.
func (d Date) MarshalText() ([]byte, error) {
	t := time.Time(d)
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalText: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(RFC3339Date))
	return t.AppendFormat(b, RFC3339Date), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in RFC 3339 format.
func (d *Date) UnmarshalText(data []byte) error {
	// Fractional seconds are handled implicitly by Parse.
	t, err := time.Parse(RFC3339Date, string(data))
	if err != nil {
		return err
	}
	*d = Date(t)
	return err
}

// Scan implements the Scanner interface.
func (t *Date) Scan(value interface{}) error {
	switch x := value.(type) {
	case nil:
		*t = Date(time.Time{})
		return nil
	case time.Time:
		*t = Date(x)
		return nil
	default:
		return fmt.Errorf("null: cannot scan type %T into types.Date: %v", value, value)
	}
}

// Value implements the driver Valuer interface.
func (t Date) Value() (driver.Value, error) {
	data, _ := t.MarshalText()
	return data, nil
}
