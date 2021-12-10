package golog

import "context"

type multiLogger struct {
	loggers []Logger
}

func (t *multiLogger) Log(ctx context.Context, level Level, kvs ...interface{}) error {
	for _, l := range t.loggers {
		if err := l.Log(ctx, level, kvs...); err != nil {
			return err
		}
	}
	return nil
}

var _ Logger = (*multiLogger)(nil)

// MultiLogger creates a logger that duplicates its logs to all the
// provided loggers, similar to the Unix tee(1) command.
//
// Each log is logged to each listed logger, one at a time.
// If a listed logger returns an error, that overall log operation
// stops and returns the error; it does not continue down the list.
func MultiLogger(loggers ...Logger) Logger {
	allLoggers := make([]Logger, 0, len(loggers))
	for _, l := range loggers {
		if ml, ok := l.(*multiLogger); ok {
			allLoggers = append(allLoggers, ml.loggers...)
		} else {
			allLoggers = append(allLoggers, l)
		}
	}
	return &multiLogger{allLoggers}
}
