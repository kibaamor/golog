package golog

import (
	"bytes"
	"context"
	"testing"
)

// Test that stdLogger properly record logs.
func TestStdLogger(t *testing.T) {
	level := LevelInfo
	tests := []struct {
		name string
		kvs  []interface{}
		want string
	}{
		{
			name: "Empty key value",
			kvs:  nil,
			want: "",
		},
		{
			name: "One key value",
			kvs:  []interface{}{"key1", "value1"},
			want: `INFO, "key1": "value1"` + "\n",
		},
		{
			name: "Two key values",
			kvs:  []interface{}{"k1", 1, "k2", 2},
			want: `INFO, "k1": "1", "k2": "2"` + "\n",
		},
		{
			name: "One key without value",
			kvs:  []interface{}{"k1"},
			want: `INFO, "k1": "KEY VALUES UNPAIRED"` + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			log := NewStdLogger(&buf)

			if err := log.Log(context.Background(), level, tt.kvs...); err != nil {
				t.Error(err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("buf.String() = %q want %q", got, tt.want)
			}
		})
	}
}
