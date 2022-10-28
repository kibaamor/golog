package golog

import (
	"bytes"
	"context"
	"testing"
)

func TestTermLogger(t *testing.T) {
	t.Parallel()

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
			want: "[INFO] key1:value1\n",
		},
		{
			name: "Two key values",
			kvs:  []interface{}{"k1", 1, "k2", 2},
			want: "[INFO] k1:1 k2:2\n",
		},
		{
			name: "One key without value",
			kvs:  []interface{}{"k1"},
			want: "[INFO] k1:KEY VALUES UNPAIRED\n",
		},
		{
			name: "Timestamp key",
			kvs:  []interface{}{DefaultTimestampKeyName, 0},
			want: "[0][INFO]\n",
		},
		{
			name: "Timestamp key with key value",
			kvs:  []interface{}{DefaultTimestampKeyName, 0, "k1", "v1"},
			want: "[0][INFO] k1:v1\n",
		},
		{
			name: "Caller key",
			kvs:  []interface{}{DefaultCallerKeyName, "caller"},
			want: "[caller][INFO]\n",
		},
		{
			name: "Caller key with key value",
			kvs:  []interface{}{DefaultCallerKeyName, "caller", "k1", "v1"},
			want: "[caller][INFO] k1:v1\n",
		},
		{
			name: "Message key",
			kvs:  []interface{}{DefaultMsgKey, "msg value"},
			want: "[INFO] msg value\n",
		},
		{
			name: "Message key with key value",
			kvs:  []interface{}{DefaultMsgKey, "msg value", "k1", "v1"},
			want: "[INFO] msg value k1:v1\n",
		},
		{
			name: "Timestamp, Caller, Message key",
			kvs:  []interface{}{DefaultTimestampKeyName, 0, DefaultCallerKeyName, "caller", DefaultMsgKey, "msg value"},
			want: "[0][caller][INFO] msg value\n",
		},
		{
			name: "Timestamp, Caller, Message key with key value",
			kvs:  []interface{}{DefaultTimestampKeyName, 0, DefaultCallerKeyName, "caller", DefaultMsgKey, "msg value", "k1", "v1"},
			want: "[0][caller][INFO] msg value k1:v1\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log := NewTermLogger(&buf, false)

			log.Log(context.Background(), level, tt.kvs...)
			if got := buf.String(); got != tt.want {
				t.Errorf("buf.String() = %q want %q", got, tt.want)
			}
		})
	}
}
