package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const strInvalidID = "xid: invalid ID"

type IDParts struct {
	id        ID
	timestamp int64
	counter   uint64
}

var IDs = []IDParts{
	IDParts{
		ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0x00, 0x10, 0xe4, 0x28, 0x41, 0x2d, 0xc9},
		5586963120321986560,
		0x10e428412dc9,
	},
	IDParts{
		ID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		0,
		0x00,
	},
	IDParts{
		ID{0x00, 0x00, 0x00, 0x00, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0x00, 0x00, 0x01},
		2864381952,
		0xccddee000001,
	},
}

func TestIDPartsExtraction(t *testing.T) {
	for i, v := range IDs {
		assert.Equal(t, v.id.Time().UnixNano(), v.timestamp, "#%d timestamp", i)
		assert.Equal(t, v.id.Counter(), v.counter, "#%d counter", i)
	}
}

func TestNew(t *testing.T) {
	// Generate 10 ids
	ids := make([]ID, 10)
	for i := 0; i < 10; i++ {
		ids[i] = NewID()
	}
	for i := 1; i < 10; i++ {
		prevID := ids[i-1]
		id := ids[i]
		// Test for uniqueness among all other 9 generated ids
		for j, tid := range ids {
			if j != i {
				assert.NotEqual(t, id, tid, "Generated ID is not unique")
			}
		}
		// Check that timestamp was incremented and is within 30 seconds of the previous one
		secs := id.Time().Sub(prevID.Time()).Seconds()
		assert.Equal(t, (secs >= 0 && secs <= 30), true, "Wrong timestamp in generated ID")
	}
}

func TestNextAfter(t *testing.T) {
	id := ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2d, 0xc9}
	assert.Equal(t, ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2d, 0xca}, id.NextAfter())
	id = ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2d, 0xff}
	assert.Equal(t, ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2e, 0x00}, id.NextAfter())
	id = ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0xff, 0xff}
	assert.Equal(t, ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x42, 0x00, 0x00}, id.NextAfter())
	id = ID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	assert.Equal(t, ID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, id.NextAfter())
}

func TestIDString(t *testing.T) {
	id := ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2d, 0xc9}
	assert.Equal(t, "9m4e2mr0ui3e8a215n4g", id.String())
}

func TestIDFromString(t *testing.T) {
	id, err := IDFromString("9m4e2mr0ui3e8a215n4g")
	assert.NoError(t, err)
	assert.Equal(t, ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2d, 0xc9}, id)
}

func TestIDFromStringInvalid(t *testing.T) {
	id, err := IDFromString("invalid")
	assert.EqualError(t, err, strInvalidID)
	assert.Equal(t, ID{}, id)
}

type jsonType struct {
	ID  *ID
	Str string
}

func TestIDJSONMarshaling(t *testing.T) {
	id := ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2d, 0xc9}
	v := jsonType{ID: &id, Str: "test"}
	data, err := json.Marshal(&v)
	assert.NoError(t, err)
	assert.Equal(t, `{"ID":"9m4e2mr0ui3e8a215n4g","Str":"test"}`, string(data))
}

func TestIDJSONUnmarshaling(t *testing.T) {
	data := []byte(`{"ID":"9m4e2mr0ui3e8a215n4g","Str":"test"}`)
	v := jsonType{}
	err := json.Unmarshal(data, &v)
	assert.NoError(t, err)
	assert.Equal(t, ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2d, 0xc9}, *v.ID)
}

func TestIDJSONUnmarshalingError(t *testing.T) {
	v := jsonType{}
	err := json.Unmarshal([]byte(`{"ID":"9M4E2MR0UI3E8A215N4G"}`), &v)
	assert.EqualError(t, err, strInvalidID)
	err = json.Unmarshal([]byte(`{"ID":"TYjhW2D0huQoQS"}`), &v)
	assert.EqualError(t, err, strInvalidID)
	err = json.Unmarshal([]byte(`{"ID":"TYjhW2D0huQoQS3kdk"}`), &v)
	assert.EqualError(t, err, strInvalidID)
}

func TestIDDriverValue(t *testing.T) {
	id := ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2d, 0xc9}
	data, err := id.Value()
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2d, 0xc9}, data)
}

func TestIDDriverScan(t *testing.T) {
	id := ID{}
	err := id.Scan("9m4e2mr0ui3e8a215n4g")
	assert.NoError(t, err)
	assert.Equal(t, ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2d, 0xc9}, id)
}

func TestIDDriverScanError(t *testing.T) {
	id := ID{}
	err := id.Scan(0)
	assert.EqualError(t, err, "xid: scanning unsupported type: int")
	err = id.Scan("0")
	assert.EqualError(t, err, strInvalidID)
	err = id.Scan([]byte{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2d})
	assert.EqualError(t, err, "xid: scanning byte slice invalid length: 11")
}

func TestIDDriverScanByteFromDatabase(t *testing.T) {
	id := ID{}
	bs := []byte{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2d, 0xc9}
	err := id.Scan(bs)
	assert.NoError(t, err)
	assert.Equal(t, ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2d, 0xc9}, id)
}

func BenchmarkNew(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = NewID()
		}
	})
}

func BenchmarkNewString(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = NewID().String()
		}
	})
}

func BenchmarkIDFromString(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = IDFromString("9m4e2mr0ui3e8a215n4g")
		}
	})
}

// func BenchmarkUUIDv1(b *testing.B) {
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			_ = uuid.NewV1().String()
// 		}
// 	})
// }

// func BenchmarkUUIDv4(b *testing.B) {
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			_ = uuid.NewV4().String()
// 		}
// 	})
// }
