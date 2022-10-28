package golog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"sync"
)

type stdLogger struct {
	log  *log.Logger
	pool *sync.Pool
}

// NewStdLogger new a standard logger with writer.
func NewStdLogger(w io.Writer) Logger {
	return &stdLogger{
		log: log.New(w, "", 0),
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// Log write the kv pairs log.
func (l *stdLogger) Log(ctx context.Context, level Level, kvs ...interface{}) {
	if len(kvs) == 0 {
		return
	}

	if (len(kvs) & 1) == 1 {
		kvs = append(kvs, "KEY VALUES UNPAIRED")
	}

	buf := l.pool.Get().(*bytes.Buffer)
	_, _ = buf.WriteString(level.String())
	for i := 0; i < len(kvs); i += 2 {
		_, _ = fmt.Fprintf(buf, `, "%v": "%v"`, kvs[i], kvs[i+1])
	}
	_ = l.log.Output(0, buf.String())
	buf.Reset()
	l.pool.Put(buf)
}
