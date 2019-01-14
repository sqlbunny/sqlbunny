package bunny

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

func IsErrNoRows(err error) bool {
	return errors.Cause(err) == sql.ErrNoRows
}

var ErrMultipleRows = errors.New("sqlbunny: multiple rows in result set")

func IsErrMultipleRows(err error) bool {
	return errors.Cause(err) == ErrMultipleRows
}

type InvalidEnumError struct {
	Value []byte
	Type  string
}

func (e *InvalidEnumError) Error() string {
	return fmt.Sprintf("Invalid %s '%s'", e.Type, e.Value)
}
