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

			if err := log.Log(context.Background(), tt.l, kvs...); err != nil {
				t.Error(err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("got %v, want: %v", got, tt.want)
			}
		})
	}
}

// Test that HandlerTimestamp properly append timestamp information into log.
func TestHandlerTimestamp(t *testing.T) {
	now := time.Now()
	nowFunc := func() time.Time {
		return now
	}
	keyName := DefaultTimestampKeyName
	valueFormat := DefaultTimestampValueFormat

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

			if err := log.Log(context.Background(), tt.l, tt.kvs...); err != nil {
				t.Error(err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("got %v, want: %v", got, tt.want)
			}
		})
	}
}

// Test that HandlerTimestamp properly append timestamp information into log.
func TestHandlerDefaultCaller(t *testing.T) {
	var buf bytes.Buffer
	log := NewStdLogger(&buf)
	log = WithHandler(log, HandlerDefaultCaller)

	if err := log.Log(context.Background(), LevelInfo, "k1", "v1"); err != nil {
		t.Error(err)
	}
	want := `INFO, "k1": "v1", "caller": "log_test.go:125"` + "\n"
	if got := buf.String(); got != want {
		t.Errorf("got %v, want: %v", got, want)
	}
}
