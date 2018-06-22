package boil

import (
	"database/sql"

	"github.com/pkg/errors"
)

func IsErrNoRows(err error) bool {
	return errors.Cause(err) == sql.ErrNoRows
}
