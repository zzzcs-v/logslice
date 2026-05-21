package query

import (
	"encoding/json"
	"testing"
)

func rawMsg(t *testing.T, s string) json.RawMessage {
	t.Helper()
	return json.RawMessage(s)
}

func TestSortEmptySlice(t *testing.T) {
	result := SortEntries(nil, SortOptions{Field: "level"})
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestSortNoField(t *testing.T) {
	entries := []json.RawMessage{
		rawMsg(t, `{"level":"error"}`),
		rawMsg(t, `{"level":"info"}`),
	}
	result := SortEntries(entries, SortOptions{})
	if string(result[0]) != `{"level":"error"}` {
		t.Errorf("expected original order preserved")
	}
}

func TestSortAscending(t *testing.T) {
	entries := []json.RawMessage{
		rawMsg(t, `{"level":"warn"}`),
		rawMsg(t, `{"level":"debug"}`),
		rawMsg(t, `{"level":"info"}`),
	}
	result := SortEntries(entries, SortOptions{Field: "level", Order: Ascending})
	expected := []string{`{"level":"debug"}`, `{"level":"info"}`, `{"level":"warn"}`}
	for i, e := range expected {
		if string(result[i]) != e {
			t.Errorf("index %d: expected %s, got %s", i, e, string(result[i]))
		}
	}
}

func TestSortDescending(t *testing.T) {
	entries := []json.RawMessage{
		rawMsg(t, `{"level":"debug"}`),
		rawMsg(t, `{"level":"warn"}`),
		rawMsg(t, `{"level":"info"}`),
	}
	result := SortEntries(entries, SortOptions{Field: "level", Order: Descending})
	expected := []string{`{"level":"warn"}`, `{"level":"info"}`, `{"level":"debug"}`}
	for i, e := range expected {
		if string(result[i]) != e {
			t.Errorf("index %d: expected %s, got %s", i, e, string(result[i]))
		}
	}
}

func TestSortMissingFieldAtEnd(t *testing.T) {
	entries := []json.RawMessage{
		rawMsg(t, `{"level":"info"}`),
		rawMsg(t, `{"msg":"no level"}`),
		rawMsg(t, `{"level":"debug"}`),
	}
	result := SortEntries(entries, SortOptions{Field: "level", Order: Ascending})
	if string(result[2]) != `{"msg":"no level"}` {
		t.Errorf("expected missing-field entry last, got %s", string(result[2]))
	}
}

func TestSortNumericField(t *testing.T) {
	entries := []json.RawMessage{
		rawMsg(t, `{"code":30}`),
		rawMsg(t, `{"code":10}`),
		rawMsg(t, `{"code":20}`),
	}
	result := SortEntries(entries, SortOptions{Field: "code", Order: Ascending})
	expected := []string{`{"code":10}`, `{"code":20}`, `{"code":30}`}
	for i, e := range expected {
		if string(result[i]) != e {
			t.Errorf("index %d: expected %s, got %s", i, e, string(result[i]))
		}
	}
}
