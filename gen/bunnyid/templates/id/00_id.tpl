{{- $modelName := .IDType.Name | titleCase -}}
{{- $modelNameCamel := .IDType.Name | camelCase -}}

import (
	"crypto/rand"
	"database/sql/driver"
    "encoding/binary"
    "encoding/json"
	"encoding/hex"
	"fmt"
	"time"
    "bytes"
    "strings"

    "github.com/pkg/errors"
    "github.com/sqlbunny/sqlbunny/runtime/bunny"
    "github.com/sqlbunny/sqlbunny/runtime/bunnyid"
)

const {{$modelNameCamel}}PrefixLength = {{ len .IDType.Prefix }} + 1
var {{$modelNameCamel}}Prefix = []byte("{{.IDType.Prefix}}_")

// ID represents a unique {{$modelName}} id
type {{$modelName}} struct{
    raw bunnyid.Raw
}

func (id {{$modelName}}) Raw() bunnyid.Raw {
    return id.raw
}

// New{{$modelName}} generates a unique {{$modelName}}
func New{{$modelName}}() {{$modelName}} {
    return {{$modelName}}FromTime(time.Now())
}

// NextAfter returns the next sequential ID after prev.
func (id {{$modelName}}) NextAfter() {{$modelName}} {
	for i := bunnyid.RawLen - 1; i >= 0; i-- {
		id.raw[i]++
		if id.raw[i] != 0 {
			break
		}
	}
	return id
}

// After returns true if this ID is after the given ID in chronological order.
func (id {{$modelName}}) After(other {{$modelName}}) bool {
	for i := 0; i < bunnyid.RawLen; i++ {
        if id.raw[i] > other.raw[i] {
            return true
        }
        if id.raw[i] < other.raw[i] {
            return false
        }
	}
	return false
}

// Before returns true if this ID is before the given ID in chronological order.
func (id {{$modelName}}) Before(other {{$modelName}}) bool {
	for i := 0; i < bunnyid.RawLen; i++ {
		if id.raw[i] < other.raw[i] {
            return true
        }
        if id.raw[i] > other.raw[i] {
            return false
        }
	}
	return false
}

// {{$modelName}}FromString reads an ID from its string representation
func {{$modelName}}FromString(id string) ({{$modelName}}, error) {
	i := &{{$modelName}}{}
	err := i.UnmarshalText([]byte(id))
	return *i, err
}

// {{$modelName}}FromTime creates an ID from the given time.
func {{$modelName}}FromTime(t time.Time) {{$modelName}} {
    var id {{$modelName}}
	binary.BigEndian.PutUint64(id.raw[:], uint64(t.UnixNano()))
	if _, err := rand.Read(id.raw[6:12]); err != nil {
		panic(errors.Errorf("cannot generate random number: %v;", err))
	}
	return id
}

// String returns a base32 hex lowercased with no padding representation of the id (char set is 0-9, a-v).
func (id {{$modelName}}) String() string {
	text := make([]byte, {{$modelNameCamel}}PrefixLength + bunnyid.EncodedLen)
    copy(text, {{$modelNameCamel}}Prefix)
	id.raw.Encode(text[{{$modelNameCamel}}PrefixLength:])
	return string(text)
}

// MarshalText implements encoding/text TextMarshaler interface
func (id {{$modelName}}) MarshalText() ([]byte, error) {
	text := make([]byte, {{$modelNameCamel}}PrefixLength + bunnyid.EncodedLen)
    copy(text, {{$modelNameCamel}}Prefix)
    id.raw.Encode(text[{{$modelNameCamel}}PrefixLength:])
	return text, nil
}

func (id {{$modelName}}) MarshalJSON() ([]byte, error) {
	text := make([]byte, {{$modelNameCamel}}PrefixLength + bunnyid.EncodedLen+2)
	text[0] = '"'
    copy(text[1:], {{$modelNameCamel}}Prefix)
    id.raw.Encode(text[1+{{$modelNameCamel}}PrefixLength:])
    text[len(text)-1] = '"'
	return text, nil
}

// UnmarshalText implements encoding/text TextUnmarshaler interface
func (id *{{$modelName}}) UnmarshalText(text []byte) error {
    if len(text) < {{$modelNameCamel}}PrefixLength {
        return &bunnyid.InvalidIDError{Value: text, Type: "{{.IDType.Name}}"}
	}
    if !bytes.Equal(text[:{{$modelNameCamel}}PrefixLength], {{$modelNameCamel}}Prefix) {
		parts := strings.Split(string(text), "_")
		if idType, ok := idPrefixes[parts[0]]; ok {
            return &bunnyid.InvalidIDError{Value: text, Type: "{{.IDType.Name}}", DetectedType: idType}
		}
        return &bunnyid.InvalidIDError{Value: text, Type: "{{.IDType.Name}}"}
	}
    if len(text) != {{$modelNameCamel}}PrefixLength + bunnyid.EncodedLen {
        return &bunnyid.InvalidIDError{Value: text, Type: "{{.IDType.Name}}"}
	}
    text = text[{{$modelNameCamel}}PrefixLength:]
    if !id.raw.Decode(text) {
        return &bunnyid.InvalidIDError{Value: text, Type: "{{.IDType.Name}}"}
    }
	return nil
}

// Time returns the timestamp part of the id.
// It's a runtime error to call this method with an invalid id.
func (id {{$modelName}}) Time() time.Time {
	// First 6 bytes of ObjectId is 64-bit big-endian nanos from epoch.
	var nowBytes [8]byte
	copy(nowBytes[0:6], id.raw[0:6])
	nanos := int64(binary.BigEndian.Uint64(nowBytes[:]))
	return time.Unix(0, nanos).UTC()
}

// Counter returns the random value part of the id.
// It's a runtime error to call this method with an invalid id.
func (id {{$modelName}}) Counter() uint64 {
	b := id.raw[6:]
	// Counter is stored as big-endian 6-byte value
	return uint64(uint64(b[0])<<40 | uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[4])<<8 | uint64(b[5]))
}

// Value implements the driver.Valuer interface.
func (id {{$modelName}}) Value() (driver.Value, error) {
	return id.raw[:], nil
}

// Scan implements the sql.Scanner interface.
func (id *{{$modelName}}) Scan(value interface{}) (err error) {
	switch val := value.(type) {
	case string:
		return id.UnmarshalText([]byte(val))
	case []byte:
		if len(val) != 12 {
			return errors.Errorf("xid: scanning byte slice invalid length: %d", len(val))
		}
		copy(id.raw[:], val[:])
		return nil
	case nil:
		*id = {{$modelName}}{}
		return nil
	default:
		return errors.Errorf("xid: scanning unsupported type: %T", value)
	}
}
