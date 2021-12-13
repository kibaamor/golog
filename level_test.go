package golog

import (
	"testing"
)

// Test that String() properly convert level to string.
func TestLevelString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    Level
		want string
	}{
		{
			name: "DEBUG",
			l:    LevelDebug,
			want: "DEBUG",
		},
		{
			name: "INFO",
			l:    LevelInfo,
			want: "INFO",
		},
		{
			name: "WARN",
			l:    LevelWarn,
			want: "WARN",
		},
		{
			name: "ERROR",
			l:    LevelError,
			want: "ERROR",
		},
		{
			name: "FATAL",
			l:    LevelFatal,
			want: "FATAL",
		},
		{
			name: "other",
			l:    10,
			want: "10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("tt.l.String() = %q want %q", got, tt.want)
			}
		})
	}
}

// Test that ParseLevel() properly parse level from string.
func TestParseLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    string
		want Level
	}{
		{
			name: "DEBUG",
			s:    "debug",
			want: LevelDebug,
		},
		{
			name: "INFO",
			s:    "inFo",
			want: LevelInfo,
		},
		{
			name: "WARN",
			s:    "WarN",
			want: LevelWarn,
		},
		{
			name: "ERROR",
			s:    "erRor",
			want: LevelError,
		},
		{
			name: "FATAL",
			s:    "FATAL",
			want: LevelFatal,
		},
		{
			name: "other",
			s:    "other",
			want: LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseLevel(tt.s); got != tt.want {
				t.Errorf("ParseLevel(tt.s) = %v want %v", got, tt.want)
			}
		})
	}
}
