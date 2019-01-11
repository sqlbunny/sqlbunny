import (
    "github.com/kernelpayments/sqlbunny/runtime/strmangle"
	"github.com/kernelpayments/sqlbunny/runtime/queries"
	"github.com/pkg/errors"
)

// M type is for providing fields and field values to UpdateAll.
type M map[string]interface{}

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
