package golog

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	// DefaultTimestampKeyName is default log key for timestamp.
	DefaultTimestampKeyName = "ts"
	// DefaultTimestampFormat is default format for timestamp.
	DefaultTimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	// DefaultTimestampNowFunc is default function for get current timestamp.
	DefaultTimestampNowFunc = time.Now
	// DefaultCallerKeyName is default log key for caller information.
	DefaultCallerKeyName = "caller"
	// DefaultCallerDepth is default depth to skip while get stack information.
	DefaultCallerDepth = 2
	// DefaultCallerWithFullPath control whether if recording the full path of log source file.
	DefaultCallerWithFullPath = false

	// HandlerDefaultTimestamp is default timestamp handler with default settings.
	HandlerDefaultTimestamp = HandlerTimestamp(DefaultTimestampKeyName, DefaultTimestampFormat, DefaultTimestampNowFunc)
	// HandlerDefaultCaller is default caller handler with default settings.
	HandlerDefaultCaller = HandlerCaller(DefaultCallerKeyName, DefaultCallerDepth, DefaultCallerWithFullPath)
)

// Filter discard log with condition.
type Filter func(ctx context.Context, level Level, kvs []interface{}) bool

// Handler modify log with anything.
type Handler func(ctx context.Context, level Level, kvs []interface{}) []interface{}

type decoratedLogger struct {
	logger  Logger
	filter  []Filter
	handler []Handler
}

func (l *decoratedLogger) Log(ctx context.Context, level Level, kvs ...interface{}) {
	for _, f := range l.filter {
		if f(ctx, level, kvs) {
			return
		}
	}
	for _, f := range l.handler {
		kvs = f(ctx, level, kvs)
	}
	l.logger.Log(ctx, level, kvs...)
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
func HandlerCaller(keyName string, depth int, withFullPath bool) Handler {
	return func(ctx context.Context, level Level, kvs []interface{}) []interface{} {
		_, file, line, _ := runtime.Caller(depth)

		// skip caller in file helper.go
		for strings.HasSuffix(file, "/helper.go") {
			depth++
			_, file, line, _ = runtime.Caller(depth)
		}

		if !withFullPath {
			index := strings.LastIndexByte(file, '/')
			file = file[index+1:]
		}
		value := file + ":" + strconv.Itoa(line)
		return append(kvs, keyName, value)
	}
}
