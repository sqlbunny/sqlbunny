package gen

import (
	"container/heap"
	"errors"
	"fmt"
	"strings"

	"github.com/sqlbunny/sqlbunny/schema"
)

type task struct {
	order int
	fn    func()
}
type taskQueue []task

func (q taskQueue) Len() int           { return len(q) }
func (q taskQueue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q taskQueue) Less(i, j int) bool { return q[i].order < q[j].order }

func (q *taskQueue) Push(x interface{}) {
	*q = append(*q, x.(task))
}

func (q *taskQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	*q = old[0 : n-1]
	return item
}

type Context struct {
	Schema *schema.Schema

	errors []error
	queue  taskQueue
}

func (ctx *Context) AddError(message string, args ...interface{}) {
	ctx.errors = append(ctx.errors, fmt.Errorf(message, args...))
}

func (ctx *Context) Enqueue(order int, fn func()) {
	heap.Push(&ctx.queue, task{order, fn})
}
func (ctx *Context) Run() {
	for len(ctx.queue) != 0 {
		t := heap.Pop(&ctx.queue).(task)
		t.fn()
	}
}

func (ctx *Context) GetType(name string, where string) schema.Type {
	res, ok := ctx.Schema.Types[name]
	if !ok {
		ctx.AddError("%s references unknown type '%s'", where, name)
	}
	return res
}

func (ctx *Context) Error() error {
	if len(ctx.errors) != 0 {
		var b strings.Builder
		fmt.Fprintf(&b, "%d errors found:\n", len(ctx.errors))
		for _, e := range ctx.errors {
			b.WriteString(e.Error())
			b.WriteRune('\n')
		}
		return errors.New(b.String())
	}
	return nil
}
