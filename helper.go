package golog

import (
	"context"
	"fmt"
	"os"
)

var (
	// DefaultMsgKey is default message key for logging.
	DefaultMsgKey = "msg"

	// DefaultMsgContext is default message Context for logging.
	DefaultMsgContext = context.Background()
)

// Option is Helper option.
type Option func(h *Helper)

// MessageKey set default log message key.
func MessageKey(key string) Option {
	return func(h *Helper) {
		h.key = key
	}
}

// Helper is a logger helper.
type Helper struct {
	logger Logger
	key    string
}

// NewHelper new a logger helper.
func NewHelper(logger Logger, opts ...Option) *Helper {
	helper := &Helper{
		logger: logger,
		key:    DefaultMsgKey,
	}
	for _, o := range opts {
		o(helper)
	}
	return helper
}

// Logger get inner logger.
func (h *Helper) Logger() Logger {
	return h.logger
}

// WithKey create a logger helper with new message key from an exist Helper.
func (h *Helper) WithKey(key string) *Helper {
	return &Helper{
		logger: h.logger,
		key:    key,
	}
}

// Log log a message.
func (h *Helper) Log(level Level, kvs ...interface{}) {
	h.logger.Log(level, kvs...)
}

// Debug logs a message at debug level.
func (h *Helper) Debug(a ...interface{}) {
	h.Log(LevelDebug, h.key, fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func (h *Helper) Debugf(format string, a ...interface{}) {
	h.Log(LevelDebug, h.key, fmt.Sprintf(format, a...))
}

// Info logs a message at info level.
func (h *Helper) Info(a ...interface{}) {
	h.Log(LevelInfo, h.key, fmt.Sprint(a...))
}

// Infof logs a message at info level.
func (h *Helper) Infof(format string, a ...interface{}) {
	h.Log(LevelInfo, h.key, fmt.Sprintf(format, a...))
}

// Warn logs a message at warn level.
func (h *Helper) Warn(a ...interface{}) {
	h.Log(LevelWarn, h.key, fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func (h *Helper) Warnf(format string, a ...interface{}) {
	h.Log(LevelWarn, h.key, fmt.Sprintf(format, a...))
}

// Error logs a message at error level.
func (h *Helper) Error(a ...interface{}) {
	h.Log(LevelError, h.key, fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func (h *Helper) Errorf(format string, a ...interface{}) {
	h.Log(LevelError, h.key, fmt.Sprintf(format, a...))
}

// Fatal logs a message at fatal level.
func (h *Helper) Fatal(a ...interface{}) {
	h.Log(LevelFatal, h.key, fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func (h *Helper) Fatalf(format string, a ...interface{}) {
	h.Log(LevelFatal, h.key, fmt.Sprintf(format, a...))
	os.Exit(1)
}
