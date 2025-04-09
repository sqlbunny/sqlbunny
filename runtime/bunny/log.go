package bunny

import (
	"context"
	"time"
)

type QueryLogInfo struct {
	Query    string
	Args     []any
	Duration time.Duration
	Err      error
}

type BeginLogInfo struct {
	ReadOnly bool
}

type CommitLogInfo struct {
	Duration time.Duration
}

type RollbackLogInfo struct {
	Duration time.Duration
	Err      error
}

type Logger interface {
	LogQuery(ctx context.Context, info QueryLogInfo)
	LogBegin(ctx context.Context, info BeginLogInfo) context.Context
	LogCommit(ctx context.Context, info CommitLogInfo)
	LogRollback(ctx context.Context, info RollbackLogInfo)
}

var logger Logger = &dummyLogger{}

func SetLogger(l Logger) {
	logger = l
}

type dummyLogger struct{}

func (l *dummyLogger) LogQuery(ctx context.Context, info QueryLogInfo)                 {}
func (l *dummyLogger) LogBegin(ctx context.Context, info BeginLogInfo) context.Context { return ctx }
func (l *dummyLogger) LogCommit(ctx context.Context, info CommitLogInfo)               {}
func (l *dummyLogger) LogRollback(ctx context.Context, info RollbackLogInfo)           {}
