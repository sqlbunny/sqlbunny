package null

import (
	"encoding/json"
	"testing"

	"github.com/volatiletech/sqlboiler/types"
)

var (
	IDVal  = types.ID{0x4d, 0x88, 0xe1, 0x5b, 0x60, 0xf4, 0x86, 0xe4, 0x28, 0x41, 0x2d, 0xc9}
	IDZero = types.ID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

func TestIDFrom(t *testing.T) {
	i := IDFrom(IDVal)
	assertID(t, i, "IDFrom()")

	zero := IDFrom(IDZero)
	if !zero.Valid {
		t.Error("IDFrom(0)", "is invalid, but should be valid")
	}
}

func TestIDFromPtr(t *testing.T) {
	iptr := &IDVal
	i := IDFromPtr(iptr)
	assertID(t, i, "IDFromPtr()")

	null := IDFromPtr(nil)
	assertNullID(t, null, "IDFromPtr(nil)")
}

func TestUnmarshalID(t *testing.T) {
	var i ID
	err := json.Unmarshal([]byte(`"9m4e2mr0ui3e8a215n4g"`), &i)
	maybePanic(err)
	assertID(t, i, "ID json")

	var null ID
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullID(t, null, "null json")

	var badType ID
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullID(t, badType, "wrong type json")

	var invalid ID
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullID(t, invalid, "invalid json")
}

func TestTextUnmarshalID(t *testing.T) {
	var i ID
	err := i.UnmarshalText([]byte("9m4e2mr0ui3e8a215n4g"))
	maybePanic(err)
	assertID(t, i, "UnmarshalText() ID")

	var blank ID
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullID(t, blank, "UnmarshalText() empty ID")
}

func TestMarshalID(t *testing.T) {
	i := IDFrom(IDVal)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, `"9m4e2mr0ui3e8a215n4g"`, "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewID(IDZero, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalIDText(t *testing.T) {
	i := IDFrom(IDVal)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "9m4e2mr0ui3e8a215n4g", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewID(IDZero, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestIDPointer(t *testing.T) {
	i := IDFrom(IDVal)
	ptr := i.Ptr()
	if *ptr != IDVal {
		t.Errorf("bad %s ID: %#v ≠ %d\n", "pointer", ptr, IDVal)
	}

	null := NewID(IDZero, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s ID: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestIDIsZero(t *testing.T) {
	i := IDFrom(IDVal)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewID(IDZero, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewID(IDZero, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestIDSetValid(t *testing.T) {
	change := NewID(IDZero, false)
	assertNullID(t, change, "SetValid()")
	change.SetValid(IDVal)
	assertID(t, change, "SetValid()")
}

func TestIDScan(t *testing.T) {
	var i ID
	err := i.Scan(IDVal[:])
	maybePanic(err)
	assertID(t, i, "scanned ID")

	var null ID
	err = null.Scan(nil)
	maybePanic(err)
	assertNullID(t, null, "scanned null")
}

func assertID(t *testing.T, i ID, from string) {
	if i.ID != IDVal {
		t.Errorf("bad %s ID: %d ≠ %d\n", from, i.ID, IDVal)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullID(t *testing.T, i ID, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
