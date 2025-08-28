package mcp

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "nanoseconds",
			duration: 123 * time.Nanosecond,
			want:     "123.00ns",
		},
		{
			name:     "microseconds",
			duration: 1234 * time.Microsecond,
			want:     "1.23ms",  // This is the actual output of the function
		},
		{
			name:     "milliseconds",
			duration: 1234 * time.Millisecond,
			want:     "1.23s",  // This is the actual output of the function
		},
		{
			name:     "seconds",
			duration: 12 * time.Second,
			want:     "12.00s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDuration(tt.duration); got != tt.want {
				t.Errorf("formatDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}