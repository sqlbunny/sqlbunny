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

import (
	"errors"
    "fmt"
    "strings"
    "github.com/kernelpayments/sqlbunny/bunny"
)

func IDFromString(s string) (bunny.ID, error) {
	parts := strings.Split(s, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Wrong parts count, expected 2 got %d", len(parts))
	}
	switch parts[0] {
    {{- range $t := .IDTypes }}
	case "{{$t.Prefix}}":
		return {{$t.Name | titleCase}}FromString(s)
    {{- end}}
	}
	return nil, fmt.Errorf("Unknown ID type %s", parts[0])
}


var idPrefixes = map[string]string{
    {{- range $t := .IDTypes }}
	"{{$t.Prefix}}": "{{$t.Name}}",
    {{- end}}
}
