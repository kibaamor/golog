package golog

// Logger is a logger interface.
type Logger interface {
	Log(level Level, kvs ...interface{})
}

// Discard is a Logger on which all Log calls succeed
// without doing anything.
var Discard Logger = discard{}

type discard struct{}

func (discard) Log(level Level, kvs ...interface{}) {
}
