# Query logging

You can set logging callbacks for every database action performed. You can use them to log interesting events with your preferred log library.

```go
type logger struct{}
func (*logger) LogQuery(ctx context.Context, info bunny.QueryLogInfo) {
	log.Printf("Query: %s", info.Query)
}
func (*logger) LogBegin(ctx context.Context, info bunny.BeginLogInfo) context.Context {
	log.Println("Begin")
	return ctx
}
func (*logger) LogCommit(ctx context.Context, info bunny.CommitLogInfo) {
	log.Println("Commit")
}
func (*logger) LogRollback(ctx context.Context, info bunny.RollbackLogInfo) {
	log.Println("Rollback")
}
```

And then enable it like this at the beginning of your program:

```go
bunny.SetLogger(&logger{})
```

## Transaction logging

You have the opportunity to attach your own data to the context in the `LogBegin` call. The returned context is passed to the inner transactio function.

This can be useful if you're attaching logging info to your context, so that all log entries made from within a transaction are annotated like so.