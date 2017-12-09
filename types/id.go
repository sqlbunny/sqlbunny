// Package xid is a globally unique id generator
//
//   - 6-byte value representing the seconds since the Unix epoch
//   - 6-byte random value
//
// The binary representation of the id is compatible with Mongo 12 bytes Object IDs.
// The string representation is using base32 hex (w/o padding) for better space efficiency
// when stored in that form (20 bytes). The hex variant of base32 is used to retain the
// sortable property of the id.
//
// Xid doesn't use base64 because case sensitivity and the 2 non alphanum chars may be an
// issue when transported as a string between various systems. Base36 wasn't retained either
// because 1/ it's not standard 2/ the resulting size is not predictable (not bit aligned)
// and 3/ it would not remain sortable. To validate a base32 `xid`, expect a 20 chars long,
// all lowercase sequence of `a` to `v` letters and `0` to `9` numbers (`[0-9a-v]{20}`).
//
// UUID is 16 bytes (128 bits), snowflake is 8 bytes (64 bits), xid stands in between
// with 12 bytes with a more compact string representation ready for the web and no
// required configuration or central generation server.
//
// Features:
//
//   - Size: 12 bytes (96 bits), smaller than UUID, larger than snowflake
//   - Base32 hex encoded by default (16 bytes storage when transported as printable string)
//   - Non configured, you don't need set a unique machine and/or data center id
//   - K-ordered
//   - Embedded time with 6 byte precision
//
// References:
//
//   - http://www.slideshare.net/davegardnerisme/unique-id-generation-in-distributed-systems
//   - https://en.wikipedia.org/wiki/Universally_unique_identifier
//   - https://blog.twitter.com/2010/announcing-snowflake
package types

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// Code inspired from mgo/bson ObjectId

// ID represents a unique request id
type ID [rawLen]byte

const (
	encodedLen = 20 // string encoded len
	rawLen     = 12 // binary raw len

	// encoding stores a custom version of the base32 encoding with lower case
	// letters.
	encoding = "0123456789abcdefghijklmnopqrstuv"
)

// ErrInvalidID is returned when trying to unmarshal an invalid ID
var ErrInvalidID = errors.New("xid: invalid ID")

// dec is the decoding map for base32 encoding
var dec [256]byte

func init() {
	for i := 0; i < len(dec); i++ {
		dec[i] = 0xFF
	}
	for i := 0; i < len(encoding); i++ {
		dec[encoding[i]] = byte(i)
	}
}

// NewID generates a globaly unique ID
func NewID() ID {
	var id ID
	binary.BigEndian.PutUint64(id[:], uint64(time.Now().UnixNano()))
	if _, err := rand.Read(id[6:12]); err != nil {
		panic(fmt.Errorf("xid: cannot generate random number: %v;", err))
	}
	return id
}

// NextAfter returns the next sequential ID after prev.
func (id ID) NextAfter() ID {
	for i := rawLen - 1; i >= 0; i-- {
		id[i]++
		if id[i] != 0 {
			break
		}
	}
	return id
}

// FromString reads an ID from its string representation
func IDFromString(id string) (ID, error) {
	i := &ID{}
	err := i.UnmarshalText([]byte(id))
	return *i, err
}

// String returns a base32 hex lowercased with no padding representation of the id (char set is 0-9, a-v).
func (id ID) String() string {
	text := make([]byte, encodedLen)
	id.encode(text)
	return string(text)
}

// MarshalText implements encoding/text TextMarshaler interface
func (id ID) MarshalText() ([]byte, error) {
	text := make([]byte, encodedLen)
	id.encode(text)
	return text, nil
}

func (id ID) MarshalJSON() ([]byte, error) {
	text := make([]byte, encodedLen+2)
	text[0] = '"'
	text[encodedLen+1] = '"'
	id.encode(text[1:])
	return text, nil
}

// encode by unrolling the stdlib base32 algorithm + removing all safe checks
func (id ID) encode(dst []byte) {
	dst[0] = encoding[id[0]>>3]
	dst[1] = encoding[(id[1]>>6)&0x1F|(id[0]<<2)&0x1F]
	dst[2] = encoding[(id[1]>>1)&0x1F]
	dst[3] = encoding[(id[2]>>4)&0x1F|(id[1]<<4)&0x1F]
	dst[4] = encoding[id[3]>>7|(id[2]<<1)&0x1F]
	dst[5] = encoding[(id[3]>>2)&0x1F]
	dst[6] = encoding[id[4]>>5|(id[3]<<3)&0x1F]
	dst[7] = encoding[id[4]&0x1F]
	dst[8] = encoding[id[5]>>3]
	dst[9] = encoding[(id[6]>>6)&0x1F|(id[5]<<2)&0x1F]
	dst[10] = encoding[(id[6]>>1)&0x1F]
	dst[11] = encoding[(id[7]>>4)&0x1F|(id[6]<<4)&0x1F]
	dst[12] = encoding[id[8]>>7|(id[7]<<1)&0x1F]
	dst[13] = encoding[(id[8]>>2)&0x1F]
	dst[14] = encoding[(id[9]>>5)|(id[8]<<3)&0x1F]
	dst[15] = encoding[id[9]&0x1F]
	dst[16] = encoding[id[10]>>3]
	dst[17] = encoding[(id[11]>>6)&0x1F|(id[10]<<2)&0x1F]
	dst[18] = encoding[(id[11]>>1)&0x1F]
	dst[19] = encoding[(id[11]<<4)&0x1F]
}

// UnmarshalText implements encoding/text TextUnmarshaler interface
func (id *ID) UnmarshalText(text []byte) error {
	if len(text) != encodedLen {
		return ErrInvalidID
	}
	for _, c := range text {
		if dec[c] == 0xFF {
			return ErrInvalidID
		}
	}
	decode(id, text)
	return nil
}

// decode by unrolling the stdlib base32 algorithm + removing all safe checks
func decode(id *ID, src []byte) {
	id[0] = dec[src[0]]<<3 | dec[src[1]]>>2
	id[1] = dec[src[1]]<<6 | dec[src[2]]<<1 | dec[src[3]]>>4
	id[2] = dec[src[3]]<<4 | dec[src[4]]>>1
	id[3] = dec[src[4]]<<7 | dec[src[5]]<<2 | dec[src[6]]>>3
	id[4] = dec[src[6]]<<5 | dec[src[7]]
	id[5] = dec[src[8]]<<3 | dec[src[9]]>>2
	id[6] = dec[src[9]]<<6 | dec[src[10]]<<1 | dec[src[11]]>>4
	id[7] = dec[src[11]]<<4 | dec[src[12]]>>1
	id[8] = dec[src[12]]<<7 | dec[src[13]]<<2 | dec[src[14]]>>3
	id[9] = dec[src[14]]<<5 | dec[src[15]]
	id[10] = dec[src[16]]<<3 | dec[src[17]]>>2
	id[11] = dec[src[17]]<<6 | dec[src[18]]<<1 | dec[src[19]]>>4
}

// Time returns the timestamp part of the id.
// It's a runtime error to call this method with an invalid id.
func (id ID) Time() time.Time {
	// First 6 bytes of ObjectId is 64-bit big-endian nanos from epoch.
	var nowBytes [8]byte
	copy(nowBytes[0:6], id[0:6])
	nanos := int64(binary.BigEndian.Uint64(nowBytes[:]))
	return time.Unix(0, nanos)
}

// Counter returns the random value part of the id.
// It's a runtime error to call this method with an invalid id.
func (id ID) Counter() uint64 {
	b := id[6:]
	// Counter is stored as big-endian 6-byte value
	return uint64(uint64(b[0])<<40 | uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[4])<<8 | uint64(b[5]))
}

// Value implements the driver.Valuer interface.
func (id ID) Value() (driver.Value, error) {
	return id[:], nil
}

// Scan implements the sql.Scanner interface.
func (id *ID) Scan(value interface{}) (err error) {
	switch val := value.(type) {
	case string:
		return id.UnmarshalText([]byte(val))
	case []byte:
		if len(val) != 12 {
			return fmt.Errorf("xid: scanning byte slice invalid length: %d", len(val))
		}
		copy(id[:], val[:])
		return nil
	default:
		return fmt.Errorf("xid: scanning unsupported type: %T", value)
	}
}
