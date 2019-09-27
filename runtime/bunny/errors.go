package bunny

import (
	"database/sql"
	"fmt"

	"github.com/sqlbunny/errors"
)

var ErrNoRows = sql.ErrNoRows

func IsErrNoRows(err error) bool {
	return errors.Is(err, ErrNoRows)
}

var ErrMultipleRows = errors.New("sqlbunny: multiple rows in result set")

func IsErrMultipleRows(err error) bool {
	return errors.Is(err, ErrMultipleRows)
}

type InvalidEnumError struct {
	Value []byte
	Type  string
}

func (e *InvalidEnumError) Error() string {
	return fmt.Sprintf("Invalid %s '%s'", e.Type, e.Value)
}
