package golog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/fatih/color"
)

// WriteFunc is log writer function.
type WriteFunc func(w io.Writer, a ...interface{})

var colorFunc map[Level]WriteFunc

func init() {
	colorFunc = map[Level]WriteFunc{
		LevelDebug: color.New(color.FgCyan).FprintlnFunc(),
		LevelInfo:  color.New(color.FgGreen).FprintlnFunc(),
		LevelWarn:  color.New(color.FgYellow).FprintlnFunc(),
		LevelError: color.New(color.FgRed).FprintlnFunc(),
		LevelFatal: color.New(color.FgRed, color.BgWhite).FprintlnFunc(),
	}
}

type termLogger struct {
	log              *log.Logger
	colorful         bool
	pool             *sync.Pool
	defaultWriteFunc WriteFunc
}

// NewTermLogger new a optimized logger for terminal with writer.
func NewTermLogger(w io.Writer, colorful bool) Logger {
	return &termLogger{
		log:      log.New(w, "", 0),
		colorful: colorful,
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		defaultWriteFunc: func(w io.Writer, a ...interface{}) {
			fmt.Fprintln(w, a...)
		},
	}
}

func extractDefaultTSCallerMsg(kvs ...interface{}) (ts, caller, msg interface{}) {
	for i := 0; i < len(kvs); i += 2 {
		k, ok := kvs[i].(string)
		if !ok {
			continue
		}

		if k == DefaultTimestampKeyName {
			ts = kvs[i+1]
		} else if k == DefaultCallerKeyName {
			caller = kvs[i+1]
		} else if k == DefaultMsgKey {
			msg = kvs[i+1]
		} else {
			continue
		}
	}
	return
}

// Log write the kv pairs log.
func (l *termLogger) Log(ctx context.Context, level Level, kvs ...interface{}) {
	if len(kvs) == 0 {
		return
	}

	if (len(kvs) & 1) == 1 {
		kvs = append(kvs, "KEY VALUES UNPAIRED")
	}

	ts, caller, msg := extractDefaultTSCallerMsg(kvs...)
	buf := l.pool.Get().(*bytes.Buffer)
	if ts != nil {
		_, _ = fmt.Fprintf(buf, "[%v]", ts)
	}
	if caller != nil {
		_, _ = fmt.Fprintf(buf, "[%v]", caller)
	}
	_, _ = fmt.Fprintf(buf, "[%v]", level)
	if msg != nil {
		_, _ = fmt.Fprintf(buf, " %v", msg)
	}

	for i := 0; i < len(kvs); i += 2 {
		if k, ok := kvs[i].(string); ok &&
			(k == DefaultTimestampKeyName || k == DefaultCallerKeyName || k == DefaultMsgKey) {
			continue
		}
		_, _ = fmt.Fprintf(buf, ` %v:%v`, kvs[i], kvs[i+1])
	}

	writeFunc := l.defaultWriteFunc
	if l.colorful {
		if fn, ok := colorFunc[level]; ok {
			writeFunc = fn
		}
	}
	writeFunc(l.log.Writer(), buf.String())

	buf.Reset()
	l.pool.Put(buf)
}
