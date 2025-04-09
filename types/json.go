package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSON is an alias for json.RawMessage, which is
// a []byte underneath.
// JSON implements Marshal and Unmarshal.
type JSON json.RawMessage

// String output your JSON.
func (j JSON) String() string {
	return string(j)
}

// Unmarshal your JSON variable into dest.
func (j JSON) Unmarshal(dest any) error {
	return json.Unmarshal(j, dest)
}

// Marshal obj into your JSON variable.
func (j *JSON) Marshal(obj any) error {
	res, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	*j = res
	return nil
}

// UnmarshalJSON sets *j to a copy of data.
func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json: unmarshal json on nil pointer to json")
	}

	*j = append((*j)[0:0], data...)
	return nil
}

// MarshalJSON returns j as the JSON encoding of j.
func (j JSON) MarshalJSON() ([]byte, error) {
	return j, nil
}

// Value returns j as a value.
// Unmarshal into RawMessage for validation.
func (j JSON) Value() (driver.Value, error) {
	var r json.RawMessage
	if err := j.Unmarshal(&r); err != nil {
		return nil, err
	}

	return []byte(r), nil
}

// Scan stores the src in *j.
func (j *JSON) Scan(src any) error {
	var source []byte

	switch src := src.(type) {
	case string:
		source = []byte(src)
	case []byte:
		source = src
	default:
		return errors.New("incompatible type for json")
	}

	*j = JSON(append((*j)[0:0], source...))

	return nil
}
