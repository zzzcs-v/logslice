package query

import (
	"testing"
	"time"
)

func TestParseTimeRangeBothEmpty(t *testing.T) {
	tr, err := ParseTimeRange("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !tr.IsZero() {
		t.Fatal("expected zero TimeRange")
	}
}

func TestParseTimeRangeFromOnly(t *testing.T) {
	tr, err := ParseTimeRange("2024-01-01", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.From == nil {
		t.Fatal("expected From to be set")
	}
	if tr.To != nil {
		t.Fatal("expected To to be nil")
	}
}

func TestParseTimeRangeRFC3339(t *testing.T) {
	tr, err := ParseTimeRange("2024-03-01T10:00:00Z", "2024-03-01T12:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.From == nil || tr.To == nil {
		t.Fatal("expected both bounds to be set")
	}
}

func TestParseTimeRangeToBeforeFrom(t *testing.T) {
	_, err := ParseTimeRange("2024-06-01", "2024-01-01")
	if err == nil {
		t.Fatal("expected error when to < from")
	}
}

func TestParseTimeRangeInvalidFrom(t *testing.T) {
	_, err := ParseTimeRange("not-a-date", "")
	if err == nil {
		t.Fatal("expected error for invalid from")
	}
}

func TestContainsWithinRange(t *testing.T) {
	tr, _ := ParseTimeRange("2024-01-01", "2024-12-31")
	mid, _ := time.Parse("2006-01-02", "2024-06-15")
	if !tr.Contains(mid) {
		t.Fatal("expected mid to be within range")
	}
}

func TestContainsBeforeRange(t *testing.T) {
	tr, _ := ParseTimeRange("2024-06-01", "2024-12-31")
	early, _ := time.Parse("2006-01-02", "2024-01-01")
	if tr.Contains(early) {
		t.Fatal("expected early to be outside range")
	}
}

func TestContainsZeroRange(t *testing.T) {
	var tr TimeRange
	if !tr.Contains(time.Now()) {
		t.Fatal("zero range should contain any time")
	}
}
