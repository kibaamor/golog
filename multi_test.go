package golog

import (
	"bytes"
	"context"
	"runtime"
	"testing"
)

// Test that MultiLogger with multiple logger.
func TestMultiLogger(t *testing.T) {
	var buf1 bytes.Buffer
	var buf2 bytes.Buffer

	slice := []Logger{NewStdLogger(&buf1), NewStdLogger(&buf2)}
	l := MultiLogger(slice...)

	l.Log(context.Background(), LevelInfo, "k1", "v1")
	want := `INFO, "k1": "v1"` + "\n"
	if got := buf1.String(); got != want {
		t.Errorf("buf1.String() = %q want %q", got, want)
	}
	if got := buf2.String(); got != want {
		t.Errorf("buf2.String() = %q want %q", got, want)
	}
}

// Test that MultiLogger copies the input slice and is insulated from future modification.
func TestMultiLoggerCopy(t *testing.T) {
	var buf bytes.Buffer

	slice := []Logger{NewStdLogger(&buf)}
	l := MultiLogger(slice...)
	slice[0] = nil

	l.Log(context.Background(), LevelInfo, "k1", "v1")
	if got, want := buf.String(), `INFO, "k1": "v1"`+"\n"; got != want {
		t.Errorf("buf.String() = %q want %q", got, want)
	}
}

// callDepth returns the logical call depth for the given PCs.
func callDepth(callers []uintptr) (depth int) {
	frames := runtime.CallersFrames(callers)
	more := true
	for more {
		_, more = frames.Next()
		depth++
	}
	return
}

// loggerFunc is an Logger implemented by the underlying func.
type loggerFunc func(ctx context.Context, level Level, kvs ...interface{})

func (f loggerFunc) Log(ctx context.Context, level Level, kvs ...interface{}) {
	f(ctx, level, kvs...)
}

// Test that MultiLogger properly flattens chained multiLoggers.
func TestMultiLoggerSingleChainFlatten(t *testing.T) {
	pc := make([]uintptr, 1000) // 1000 should fit the full stack
	n := runtime.Callers(0, pc)
	myDepth := callDepth(pc[:n])
	var logDepth int // will contain the depth from which loggerFunc.Logger was called

	l := MultiLogger(loggerFunc(func(ctx context.Context, level Level, kvs ...interface{}) {
		n := runtime.Callers(1, pc)
		logDepth += callDepth(pc[:n])
	}))

	ml := l
	// chain a bunch of multiLoggers
	for i := 0; i < 100; i++ {
		ml = MultiLogger(ml)
	}
	ml = MultiLogger(l, ml, l, ml)
	ml.Log(context.Background(), LevelInfo, "k1", "v1")

	if logDepth != 4*(myDepth+2) { // 2 should be multiLogger.Log and loggerFunc.Log
		t.Errorf("multiLogger did not flatten chained multiLoggers: expected logDepth %d, got %d",
			4*(myDepth+2), logDepth)
	}
}
