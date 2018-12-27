import (
    "github.com/kernelpayments/sqlbunny/runtime/strmangle"
	"github.com/kernelpayments/sqlbunny/runtime/queries"
	"github.com/pkg/errors"
)

// M type is for providing fields and field values to UpdateAll.
type M map[string]interface{}

// ErrSyncFail occurs during insert when the record could not be retrieved in
// order to populate default value information. This usually happens when LastInsertId
// fails or there was a primary key configuration that was not resolvable.
var ErrSyncFail = errors.New("{{.PkgName}}: failed to synchronize data after insert")

type insertCache struct {
	query        string
	retQuery     string
	valueMapping []queries.MappedField
	retMapping   []queries.MappedField
}

type updateCache struct {
	query        string
	valueMapping []queries.MappedField
}

func makeCacheKey(wl, nzDefaults []string) string {
	buf := strmangle.GetBuffer()

	for _, w := range wl {
		buf.WriteString(w)
	}
	if len(nzDefaults) != 0 {
		buf.WriteByte('.')
	}
	for _, nz := range nzDefaults {
		buf.WriteString(nz)
	}

	str := buf.String()
	strmangle.PutBuffer(buf)
	return str
}
