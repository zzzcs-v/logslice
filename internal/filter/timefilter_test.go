package filter

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/user/logslice/internal/query"
)

func makeEntry(field, ts string) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{%q:%q,"msg":"hello"}`, field, ts))
}

func TestMatchTimeRangeZero(t *testing.T) {
	var tr query.TimeRange
	raw := makeEntry("time", "2024-06-01T00:00:00Z")
	if !MatchTimeRange(raw, tr) {
		t.Fatal("zero range should match everything")
	}
}

func TestMatchTimeRangeWithin(t *testing.T) {
	tr, _ := query.ParseTimeRange("2024-01-01T00:00:00Z", "2024-12-31T23:59:59Z")
	raw := makeEntry("time", "2024-06-15T12:00:00Z")
	if !MatchTimeRange(raw, tr) {
		t.Fatal("expected entry within range to match")
	}
}

func TestMatchTimeRangeOutside(t *testing.T) {
	tr, _ := query.ParseTimeRange("2024-06-01T00:00:00Z", "2024-06-30T23:59:59Z")
	raw := makeEntry("time", "2024-01-01T00:00:00Z")
	if MatchTimeRange(raw, tr) {
		t.Fatal("expected entry outside range to not match")
	}
}

func TestMatchTimeRangeAlternateField(t *testing.T) {
	tr, _ := query.ParseTimeRange("2024-01-01T00:00:00Z", "2024-12-31T23:59:59Z")
	raw := makeEntry("timestamp", "2024-03-10T08:30:00Z")
	if !MatchTimeRange(raw, tr) {
		t.Fatal("expected timestamp field to be recognised")
	}
}

func TestMatchTimeRangeNoTimestampField(t *testing.T) {
	tr, _ := query.ParseTimeRange("2024-01-01T00:00:00Z", "2024-12-31T23:59:59Z")
	raw := json.RawMessage(`{"level":"info","msg":"no ts"}`)
	if !MatchTimeRange(raw, tr) {
		t.Fatal("entry without timestamp should pass through")
	}
}

func TestMatchTimeRangeMalformedJSON(t *testing.T) {
	tr, _ := query.ParseTimeRange("2024-01-01T00:00:00Z", "2024-12-31T23:59:59Z")
	raw := json.RawMessage(`not json`)
	if !MatchTimeRange(raw, tr) {
		t.Fatal("malformed JSON should pass through")
	}
}

func TestMatchTimeRangeBoundaryInclusive(t *testing.T) {
	boundary := "2024-06-01T00:00:00Z"
	tr, _ := query.ParseTimeRange(boundary, boundary)
	raw := makeEntry("ts", boundary)
	if !MatchTimeRange(raw, tr) {
		t.Fatal("boundary value should be inclusive")
	}
}

func TestMatchTimeRangeAtTimestamp(t *testing.T) {
	tr, _ := query.ParseTimeRange("2024-01-01T00:00:00Z", "2024-12-31T23:59:59Z")
	ts := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	raw := makeEntry("@timestamp", ts)
	if !MatchTimeRange(raw, tr) {
		t.Fatal("expected @timestamp field to be recognised")
	}
}
