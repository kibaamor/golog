package golog

import "context"

// Logger is a logger interface.
type Logger interface {
	Log(ctx context.Context, level Level, kvs ...interface{})
}

// Discard is a Logger on which all Log calls succeed
// without doing anything.
var Discard Logger = discard{}

type discard struct{}

func (discard) Log(ctx context.Context, level Level, kvs ...interface{}) {
}
