package golog

import (
	"bytes"
	"context"
	"io"
	"testing"
)

func testHelperKey(helperFunc func(logger Logger, key string) *Helper, t *testing.T) {
	tests := []struct {
		name string
		call func(helper *Helper)
		want string
	}{
		{
			name: "debug",
			call: func(helper *Helper) {
				helper.Debug(1, "2", 3)
			},
			want: `DEBUG, "log": "123"`,
		},
		{
			name: "debugf",
			call: func(helper *Helper) {
				helper.Debugf("%d %d %d", 1, 2, 3)
			},
			want: `DEBUG, "log": "1 2 3"`,
		},
		{
			name: "info",
			call: func(helper *Helper) {
				helper.Info(1, "2", 3)
			},
			want: `INFO, "log": "123"`,
		},
		{
			name: "infof",
			call: func(helper *Helper) {
				helper.Infof("%d %d %d", 1, 2, 3)
			},
			want: `INFO, "log": "1 2 3"`,
		},
		{
			name: "warn",
			call: func(helper *Helper) {
				helper.Warn(1, "2", 3)
			},
			want: `WARN, "log": "123"`,
		},
		{
			name: "warnf",
			call: func(helper *Helper) {
				helper.Warnf("%d %d %d", 1, 2, 3)
			},
			want: `WARN, "log": "1 2 3"`,
		},
		{
			name: "error",
			call: func(helper *Helper) {
				helper.Error(1, "2", 3)
			},
			want: `ERROR, "log": "123"`,
		},
		{
			name: "errorf",
			call: func(helper *Helper) {
				helper.Errorf("%d %d %d", 1, 2, 3)
			},
			want: `ERROR, "log": "1 2 3"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			logger := NewStdLogger(&buf)
			helper := helperFunc(logger, "log")

			tt.call(helper)
			want := tt.want + "\n"

			if got := buf.String(); got != want {
				t.Errorf("buf.String() = %q want %q", got, want)
			}
		})
	}
}

// Test that NewHelper properly record logs.
func TestNewHelper(t *testing.T) {
	testHelperKey(func(logger Logger, key string) *Helper {
		return NewHelper(logger, MessageKey(key))
	}, t)
}

// Test that WithKey properly record logs.
func TestHelperWithKey(t *testing.T) {
	testHelperKey(func(logger Logger, key string) *Helper {
		helper := NewHelper(logger)
		return helper.WithKey(key)
	}, t)
}

type testKey struct{}

// Test that WithContext properly record logs.
func TestHelperWithContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), testKey{}, "test value")

	var buf bytes.Buffer
	logger := NewStdLogger(&buf)
	logger = WithHandler(logger, func(ctx context.Context, level Level, kvs []interface{}) []interface{} {
		return append(kvs, "test key", ctx.Value(testKey{}).(string))
	})
	helper := NewHelper(logger, MessageContext(ctx))

	helper.Info("k")

	if got, want := buf.String(), `INFO, "msg": "k", "test key": "test value"`+"\n"; got != want {
		t.Errorf("buf.String() = %q want %q", got, want)
	}
}

func BenchmarkHelperPrint(b *testing.B) {
	log := NewHelper(NewStdLogger(io.Discard))
	for i := 0; i < b.N; i++ {
		log.Debug("test")
	}
}

func BenchmarkHelperPrintf(b *testing.B) {
	log := NewHelper(NewStdLogger(io.Discard))
	for i := 0; i < b.N; i++ {
		log.Debugf("%s", "test")
	}
}
