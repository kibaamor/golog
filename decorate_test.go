package golog

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"
)

// Test that FilterLevel properly filter level less than specific.
func TestFilterLevel(t *testing.T) {
	t.Parallel()

	filterLevel := LevelWarn
	kvs := []interface{}{"k1", "v1"}
	tests := []struct {
		name string
		l    Level
		want string
	}{
		{
			name: "DEBUG",
			l:    LevelDebug,
			want: "",
		},
		{
			name: "INFO",
			l:    LevelInfo,
			want: "",
		},
		{
			name: "WARN",
			l:    LevelWarn,
			want: `WARN, "k1": "v1"` + "\n",
		},
		{
			name: "ERROR",
			l:    LevelError,
			want: `ERROR, "k1": "v1"` + "\n",
		},
		{
			name: "FATAL",
			l:    LevelFatal,
			want: `FATAL, "k1": "v1"` + "\n",
		},
		{
			name: "other",
			l:    10,
			want: `10, "k1": "v1"` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log := NewStdLogger(&buf)
			log = WithFilter(log, FilterLevel(filterLevel))

			log.Log(context.Background(), tt.l, kvs...)
			if got := buf.String(); got != tt.want {
				t.Errorf("buf.String() = %q want = %q", got, tt.want)
			}
		})
	}
}

// Test that HandlerTimestamp properly append timestamp information into log.
func TestHandlerTimestamp(t *testing.T) {
	t.Parallel()

	now := time.Now()
	nowFunc := func() time.Time {
		return now
	}
	keyName := DefaultTimestampKeyName
	valueFormat := DefaultTimestampFormat

	tests := []struct {
		name string
		l    Level
		kvs  []interface{}
		want string
	}{
		{
			name: "Without Log",
			l:    LevelInfo,
			kvs:  nil,
			want: fmt.Sprintf(`INFO, "%s": "%s"`+"\n", keyName, now.Format(valueFormat)),
		},
		{
			name: "With 1 Log",
			l:    LevelInfo,
			kvs:  []interface{}{"k1", 1},
			want: fmt.Sprintf(`INFO, "k1": "1", "%s": "%s"`+"\n", keyName, now.Format(valueFormat)),
		},
		{
			name: "With 2 Logs",
			l:    LevelInfo,
			kvs:  []interface{}{"k1", 1, "k2", "v2"},
			want: fmt.Sprintf(`INFO, "k1": "1", "k2": "v2", "%s": "%s"`+"\n", keyName, now.Format(valueFormat)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log := NewStdLogger(&buf)
			log = WithHandler(log, HandlerTimestamp(keyName, valueFormat, nowFunc))

			log.Log(context.Background(), tt.l, tt.kvs...)
			if got := buf.String(); got != tt.want {
				t.Errorf("buf.String() = %q want = %q", got, tt.want)
			}
		})
	}
}

// Test that HandlerTimestamp properly append timestamp information into log.
func TestHandlerDefaultCaller(t *testing.T) {
	var buf bytes.Buffer
	log := NewStdLogger(&buf)
	log = WithHandler(log, HandlerDefaultCaller)

	log.Log(context.Background(), LevelInfo, "k1", "v1")
	if got, want := buf.String(), `INFO, "k1": "v1", "caller": "decorate_test.go:125"`+"\n"; got != want {
		t.Errorf("buf.String() = %q want = %q", got, want)
	}
}
