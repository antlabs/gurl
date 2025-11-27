package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type DailyCron struct {
	Second int
	Minute int
	Hour   int
}

// ParseDailyCron parses a simple 6-field cron expression: "sec min hour * * *".
// Only daily expressions with '*' for day, month and weekday are supported.
func ParseDailyCron(expr string) (*DailyCron, error) {
	parts := strings.Fields(expr)
	if len(parts) != 6 {
		return nil, fmt.Errorf("invalid cron expression, expected 6 fields: %s", expr)
	}
	if parts[3] != "*" || parts[4] != "*" || parts[5] != "*" {
		return nil, fmt.Errorf("only daily cron expressions with '*' for day, month and weekday are supported: %s", expr)
	}
	sec, err := strconv.Atoi(parts[0])
	if err != nil || sec < 0 || sec > 59 {
		return nil, fmt.Errorf("invalid seconds field in cron expression: %s", parts[0])
	}
	min, err := strconv.Atoi(parts[1])
	if err != nil || min < 0 || min > 59 {
		return nil, fmt.Errorf("invalid minutes field in cron expression: %s", parts[1])
	}
	hour, err := strconv.Atoi(parts[2])
	if err != nil || hour < 0 || hour > 23 {
		return nil, fmt.Errorf("invalid hours field in cron expression: %s", parts[2])
	}
	return &DailyCron{Second: sec, Minute: min, Hour: hour}, nil
}

// NextAfter returns the next time after t when this cron should trigger.
func (c *DailyCron) NextAfter(t time.Time) time.Time {
	y, m, d := t.Date()
	loc := t.Location()
	next := time.Date(y, m, d, c.Hour, c.Minute, c.Second, 0, loc)
	if !next.After(t) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}
