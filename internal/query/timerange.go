package query

import (
	"fmt"
	"time"
)

// TimeRange represents an inclusive time window for filtering log entries.
type TimeRange struct {
	From *time.Time
	To   *time.Time
}

// ParseTimeRange parses "from" and "to" strings into a TimeRange.
// Accepts RFC3339 or the shorthand "2006-01-02" date format.
func ParseTimeRange(from, to string) (TimeRange, error) {
	var tr TimeRange

	if from != "" {
		t, err := parseTime(from)
		if err != nil {
			return tr, fmt.Errorf("invalid from time %q: %w", from, err)
		}
		tr.From = &t
	}

	if to != "" {
		t, err := parseTime(to)
		if err != nil {
			return tr, fmt.Errorf("invalid to time %q: %w", to, err)
		}
		tr.To = &t
	}

	if tr.From != nil && tr.To != nil && tr.To.Before(*tr.From) {
		return tr, fmt.Errorf("to time must not be before from time")
	}

	return tr, nil
}

// Contains reports whether t falls within the TimeRange (inclusive).
func (tr TimeRange) Contains(t time.Time) bool {
	if tr.From != nil && t.Before(*tr.From) {
		return false
	}
	if tr.To != nil && t.After(*tr.To) {
		return false
	}
	return true
}

// IsZero reports whether the TimeRange has no bounds set.
func (tr TimeRange) IsZero() bool {
	return tr.From == nil && tr.To == nil
}

func parseTime(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised time format")
}
