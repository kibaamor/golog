package golog

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	DefaultTimestampKeyName     = "ts"
	DefaultTimestampValueFormat = "2006-01-02 15:04:05.000"
	DefaultTimestampNowFunc     = time.Now
	DefaultCallerKeyName        = "caller"
	DefaultCallerDepth          = 2
	DefaultCallerWithFullpath   = false

	HandlerDefaultTimestamp = HandlerTimestamp(DefaultTimestampKeyName, DefaultTimestampValueFormat, DefaultTimestampNowFunc)
	HandlerDefaultCaller    = HandlerCaller(DefaultCallerKeyName, DefaultCallerDepth, DefaultCallerWithFullpath)
)

// Logger is a logger interface.
type Logger interface {
	Log(ctx context.Context, level Level, kvs ...interface{}) error
}

// Discard is a Logger on which all Log calls succeed
// without doing anything.
var Discard Logger = discard{}

type discard struct{}

func (discard) Log(ctx context.Context, level Level, kvs ...interface{}) error {
	return nil
}

// Filter discard log with condition.
type Filter func(ctx context.Context, level Level, kvs []interface{}) bool

// Handler modify log with anything.
type Handler func(ctx context.Context, level Level, kvs []interface{}) []interface{}

type decoratedLogger struct {
	logger  Logger
	filter  []Filter
	handler []Handler
}

func (l *decoratedLogger) Log(ctx context.Context, level Level, kvs ...interface{}) error {
	for _, f := range l.filter {
		if f(ctx, level, kvs) {
			return nil
		}
	}
	for _, f := range l.handler {
		kvs = f(ctx, level, kvs)
	}
	return l.logger.Log(ctx, level, kvs...)
}

// WithFilter decorate logger with filters
func WithFilter(logger Logger, filter ...Filter) Logger {
	if l, ok := logger.(*decoratedLogger); ok {
		return &decoratedLogger{
			logger:  l.logger,
			filter:  append(l.filter, filter...),
			handler: l.handler,
		}
	}
	return &decoratedLogger{logger: logger, filter: filter}
}

// WithHandler decorate logger with handlers
func WithHandler(logger Logger, handler ...Handler) Logger {
	if l, ok := logger.(*decoratedLogger); ok {
		return &decoratedLogger{
			logger:  l.logger,
			filter:  l.filter,
			handler: append(l.handler, handler...),
		}
	}
	return &decoratedLogger{logger: logger, handler: handler}
}

// FilterLevel filter log level less than specific level.
func FilterLevel(l Level) Filter {
	return func(ctx context.Context, level Level, kvs []interface{}) bool {
		return level < l
	}
}

// HandlerTimestamp append timestamp information into log.
func HandlerTimestamp(keyName, valueFormat string, nowFunc func() time.Time) Handler {
	return func(ctx context.Context, level Level, kvs []interface{}) []interface{} {
		return append(kvs, keyName, nowFunc().Format(valueFormat))
	}
}

// HandlerCaller append caller information into log.
func HandlerCaller(keyName string, depth int, withFullpath bool) Handler {
	return func(ctx context.Context, level Level, kvs []interface{}) []interface{} {
		_, file, line, _ := runtime.Caller(depth)
		for strings.LastIndex(file, "/log/helper.go") > 0 {
			depth++
			_, file, line, _ = runtime.Caller(depth)
		}
		if !withFullpath {
			index := strings.LastIndexByte(file, '/')
			file = file[index+1:]
		}
		value := file + ":" + strconv.Itoa(line)
		return append(kvs, keyName, value)
	}
}
