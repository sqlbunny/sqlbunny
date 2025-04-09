import (
	"github.com/sqlbunny/errors"
    "github.com/sqlbunny/sqlbunny/runtime/strmangle"
	"github.com/sqlbunny/sqlbunny/runtime/queries"
)

// M type is for providing fields and field values to UpdateAll.
type M map[string]any

type insertCache struct {
	query        string
	valueMapping []queries.MappedField
}

type updateCache struct {
	query        string
	valueMapping []queries.MappedField
}

func makeCacheKey(wl []string) string {
	buf := strmangle.GetBuffer()

	for _, w := range wl {
		buf.WriteString(w)
		buf.WriteByte(',')
	}

	str := buf.String()
	strmangle.PutBuffer(buf)
	return str
}
