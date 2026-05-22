package query

import (
	"encoding/json"
	"testing"
)

func rawEntry(t *testing.T, s string) json.RawMessage {
	t.Helper()
	return json.RawMessage(s)
}

func TestGroupByEmptyField(t *testing.T) {
	entries := []json.RawMessage{
		rawEntry(t, `{"level":"info"}`),
		rawEntry(t, `{"level":"error"}`),
	}
	results := GroupBy(entries, "")
	if len(results) != 1 {
		t.Fatalf("expected 1 group, got %d", len(results))
	}
	if len(results[0].Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(results[0].Entries))
	}
}

func TestGroupByNoEntries(t *testing.T) {
	results := GroupBy(nil, "level")
	if len(results) != 1 || results[0].Key != "" {
		t.Errorf("expected single empty group, got %+v", results)
	}
}

func TestGroupByField(t *testing.T) {
	entries := []json.RawMessage{
		rawEntry(t, `{"level":"info","msg":"a"}`),
		rawEntry(t, `{"level":"error","msg":"b"}`),
		rawEntry(t, `{"level":"info","msg":"c"}`),
	}
	results := GroupBy(entries, "level")
	if len(results) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(results))
	}
	// sorted: error < info
	if results[0].Key != "error" {
		t.Errorf("expected first key=error, got %q", results[0].Key)
	}
	if len(results[0].Entries) != 1 {
		t.Errorf("expected 1 error entry, got %d", len(results[0].Entries))
	}
	if results[1].Key != "info" {
		t.Errorf("expected second key=info, got %q", results[1].Key)
	}
	if len(results[1].Entries) != 2 {
		t.Errorf("expected 2 info entries, got %d", len(results[1].Entries))
	}
}

func TestGroupByMissingField(t *testing.T) {
	entries := []json.RawMessage{
		rawEntry(t, `{"level":"info"}`),
		rawEntry(t, `{"msg":"no level"}`),
	}
	results := GroupBy(entries, "level")
	keys := map[string]int{}
	for _, r := range results {
		keys[r.Key] = len(r.Entries)
	}
	if keys["info"] != 1 {
		t.Errorf("expected 1 info entry")
	}
	if keys[""] != 1 {
		t.Errorf("expected 1 missing-field entry")
	}
}

func TestGroupByMalformedJSON(t *testing.T) {
	entries := []json.RawMessage{
		rawEntry(t, `{"level":"info"}`),
		rawEntry(t, `not-json`),
	}
	results := GroupBy(entries, "level")
	total := 0
	for _, r := range results {
		total += len(r.Entries)
	}
	if total != 2 {
		t.Errorf("expected all entries accounted for, got %d", total)
	}
}
