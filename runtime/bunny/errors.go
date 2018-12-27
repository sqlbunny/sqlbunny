package bunny

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func IsErrNoRows(err error) bool {
	return errors.Cause(err) == sql.ErrNoRows
}

func IsErrUniqueViolation(err error) bool {
	cause := errors.Cause(err)
	if pqErr, ok := cause.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}

type InvalidEnumError struct {
	Value []byte
	Type  string
}

func (e *InvalidEnumError) Error() string {
	return fmt.Sprintf("Invalid %s '%s'", e.Type, e.Value)
}
