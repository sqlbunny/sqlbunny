package boil

import (
	"context"
	"time"
)

type QueryLogger interface {
	LogQuery(ctx context.Context, query string, elapsed time.Duration, args ...interface{})
}

var queryLogger QueryLogger

func SetQueryLogger(ql QueryLogger) {
	queryLogger = ql
}
