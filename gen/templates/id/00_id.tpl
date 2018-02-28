{{- $modelName := .IDType.Name | titleCase -}}
{{- $modelNameCamel := .IDType.Name | camelCase -}}

import (
	"crypto/rand"
	"database/sql/driver"
    "encoding/binary"
	"encoding/json"
	"fmt"
	"time"
    "bytes"
    "strings"

    "github.com/pkg/errors"
    "github.com/KernelPay/sqlboiler/boil"
    "github.com/KernelPay/toolkit/apierrors"
)

const {{$modelNameCamel}}PrefixLength = {{ len .IDType.Prefix }} + 1
var {{$modelNameCamel}}Prefix = []byte("{{.IDType.Prefix}}_")

// ID represents a unique request id
type {{$modelName}}ID struct{
    boil.IDData
}

func (id {{$modelName}}ID) IDBytes() []byte {
    return id.IDData[:]
}

// New{{$modelName}}ID generates a globaly unique {{$modelName}}ID
func New{{$modelName}}ID() {{$modelName}}ID {
	var id {{$modelName}}ID
	binary.BigEndian.PutUint64(id.IDData[:], uint64(time.Now().UnixNano()))
	if _, err := rand.Read(id.IDData[6:12]); err != nil {
		panic(errors.Errorf("cannot generate random number: %v;", err))
	}
	return id
}

// NextAfter returns the next sequential ID after prev.
func (id {{$modelName}}ID) NextAfter() {{$modelName}}ID {
	for i := boil.IDRawLen - 1; i >= 0; i-- {
		id.IDData[i]++
		if id.IDData[i] != 0 {
			break
		}
	}
	return id
}

// FromString reads an ID from its string representation
func {{$modelName}}IDFromString(id string) ({{$modelName}}ID, error) {
	i := &{{$modelName}}ID{}
	err := i.UnmarshalText([]byte(id))
	return *i, err
}

// String returns a base32 hex lowercased with no padding representation of the id (char set is 0-9, a-v).
func (id {{$modelName}}ID) String() string {
	text := make([]byte, {{$modelNameCamel}}PrefixLength + boil.IDEncodedLen)
    copy(text, {{$modelNameCamel}}Prefix)
	id.IDData.Encode(text[{{$modelNameCamel}}PrefixLength:])
	return string(text)
}

// MarshalText implements encoding/text TextMarshaler interface
func (id {{$modelName}}ID) MarshalText() ([]byte, error) {
	text := make([]byte, {{$modelNameCamel}}PrefixLength + boil.IDEncodedLen)
    copy(text, {{$modelNameCamel}}Prefix)
    id.IDData.Encode(text[{{$modelNameCamel}}PrefixLength:])
	return text, nil
}

func (id {{$modelName}}ID) MarshalJSON() ([]byte, error) {
	text := make([]byte, {{$modelNameCamel}}PrefixLength + boil.IDEncodedLen+2)
	text[0] = '"'
    copy(text[1:], {{$modelNameCamel}}Prefix)
    id.IDData.Encode(text[1+{{$modelNameCamel}}PrefixLength:])
    text[len(text)-1] = '"'
	return text, nil
}

// UnmarshalText implements encoding/text TextUnmarshaler interface
func (id *{{$modelName}}ID) UnmarshalText(text []byte) error {
    if !bytes.Equal(text[:{{$modelNameCamel}}PrefixLength], {{$modelNameCamel}}Prefix) {
		parts := strings.Split(string(text), "_")
		if idType, ok := idPrefixes[parts[0]]; ok {
			return apierrors.New(apierrors.TypeInvalidRequest, "Invalid {{.IDType.Name}} ID '%s': You're passing a %s ID, but we need a {{.IDType.Name}} ID here.", text, idType)
		}
		return apierrors.New(apierrors.TypeInvalidRequest, "Invalid {{.IDType.Name}} ID '%s'", text)
	}
    if len(text) != {{$modelNameCamel}}PrefixLength + boil.IDEncodedLen {
		return apierrors.New(apierrors.TypeInvalidRequest, "Invalid {{.IDType.Name}} ID '%s'", text)
	}
    text = text[{{$modelNameCamel}}PrefixLength:]
    if !id.IDData.Decode(text) {
        return apierrors.New(apierrors.TypeInvalidRequest, "Invalid {{.IDType.Name}} ID '%s'", text)
    }
	return nil
}

// Time returns the timestamp part of the id.
// It's a runtime error to call this method with an invalid id.
func (id {{$modelName}}ID) Time() time.Time {
	// First 6 bytes of ObjectId is 64-bit big-endian nanos from epoch.
	var nowBytes [8]byte
	copy(nowBytes[0:6], id.IDData[0:6])
	nanos := int64(binary.BigEndian.Uint64(nowBytes[:]))
	return time.Unix(0, nanos)
}

// Counter returns the random value part of the id.
// It's a runtime error to call this method with an invalid id.
func (id {{$modelName}}ID) Counter() uint64 {
	b := id.IDData[6:]
	// Counter is stored as big-endian 6-byte value
	return uint64(uint64(b[0])<<40 | uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[4])<<8 | uint64(b[5]))
}

// Value implements the driver.Valuer interface.
func (id {{$modelName}}ID) Value() (driver.Value, error) {
	return id.IDData[:], nil
}

// Scan implements the sql.Scanner interface.
func (id *{{$modelName}}ID) Scan(value interface{}) (err error) {
	switch val := value.(type) {
	case string:
		return id.UnmarshalText([]byte(val))
	case []byte:
		if len(val) != 12 {
			return errors.Errorf("xid: scanning byte slice invalid length: %d", len(val))
		}
		copy(id.IDData[:], val[:])
		return nil
	default:
		return errors.Errorf("xid: scanning unsupported type: %T", value)
	}
}
