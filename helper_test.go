package golog

import (
	"bytes"
	"context"
	"io"
	"testing"
)

func testHelperKey(helperFunc func(logger Logger, key string) *Helper, t *testing.T) {
	var buf bytes.Buffer
	logger := NewStdLogger(&buf)
	helper := helperFunc(logger, "log")

	want := `DEBUG, "log": "123"` + "\n"
	helper.Debug(1, "2", 3)
	if buf.String() != want {
		t.Errorf("got: %v, want: %s", buf.String(), want)
	}
	buf.Reset()

	want = `DEBUG, "log": "1 2 3"` + "\n"
	helper.Debugf("%d %d %d", 1, 2, 3)
	if buf.String() != want {
		t.Errorf("got: %v, want: %s", buf.String(), want)
	}
	buf.Reset()

	want = `INFO, "log": "123"` + "\n"
	helper.Info(1, "2", 3)
	if buf.String() != want {
		t.Errorf("got: %v, want: %s", buf.String(), want)
	}
	buf.Reset()

	want = `INFO, "log": "1 2 3"` + "\n"
	helper.Infof("%d %d %d", 1, 2, 3)
	if buf.String() != want {
		t.Errorf("got: %v, want: %s", buf.String(), want)
	}
	buf.Reset()

	want = `WARN, "log": "123"` + "\n"
	helper.Warn(1, "2", 3)
	if buf.String() != want {
		t.Errorf("got: %v, want: %s", buf.String(), want)
	}
	buf.Reset()

	want = `WARN, "log": "1 2 3"` + "\n"
	helper.Warnf("%d %d %d", 1, 2, 3)
	if buf.String() != want {
		t.Errorf("got: %v, want: %s", buf.String(), want)
	}
	buf.Reset()

	want = `ERROR, "log": "123"` + "\n"
	helper.Error(1, "2", 3)
	if buf.String() != want {
		t.Errorf("got: %v, want: %s", buf.String(), want)
	}
	buf.Reset()

	want = `ERROR, "log": "1 2 3"` + "\n"
	helper.Errorf("%d %d %d", 1, 2, 3)
	if buf.String() != want {
		t.Errorf("got: %v, want: %s", buf.String(), want)
	}
	buf.Reset()
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

	want := `INFO, "msg": "k", "test key": "test value"` + "\n"
	if buf.String() != want {
		t.Errorf("got: %v, want: %s", buf.String(), want)
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
