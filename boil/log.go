package boil

import (
	"context"
	"time"
)

type QueryLogger interface {
	LogQuery(ctx context.Context, query string, duration time.Duration, args ...interface{})
	LogBegin(ctx context.Context, readOnly bool) context.Context
	LogCommit(ctx context.Context, duration time.Duration)
	LogRollback(ctx context.Context, duration time.Duration, cause error)
}

var queryLogger QueryLogger

func SetQueryLogger(ql QueryLogger) {
	queryLogger = ql
}
