package query

import (
	"encoding/json"
	"testing"
)

func rawDedupEntry(kv map[string]interface{}) json.RawMessage {
	b, _ := json.Marshal(kv)
	return json.RawMessage(b)
}

func TestDedupEmptyField(t *testing.T) {
	entries := []json.RawMessage{
		rawDedupEntry(map[string]interface{}{"level": "info"}),
		rawDedupEntry(map[string]interface{}{"level": "info"}),
	}
	result, err := DedupEntries(entries, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestDedupNoEntries(t *testing.T) {
	result, err := DedupEntries(nil, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}

func TestDedupByField(t *testing.T) {
	entries := []json.RawMessage{
		rawDedupEntry(map[string]interface{}{"level": "info", "msg": "a"}),
		rawDedupEntry(map[string]interface{}{"level": "info", "msg": "b"}),
		rawDedupEntry(map[string]interface{}{"level": "error", "msg": "c"}),
		rawDedupEntry(map[string]interface{}{"level": "info", "msg": "d"}),
	}
	result, err := DedupEntries(entries, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestDedupMissingField(t *testing.T) {
	entries := []json.RawMessage{
		rawDedupEntry(map[string]interface{}{"msg": "no level here"}),
		rawDedupEntry(map[string]interface{}{"msg": "also no level"}),
	}
	result, err := DedupEntries(entries, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Missing field entries are preserved as-is
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestDedupSkipsMalformed(t *testing.T) {
	entries := []json.RawMessage{
		json.RawMessage(`not-json`),
		rawDedupEntry(map[string]interface{}{"level": "info"}),
	}
	result, err := DedupEntries(entries, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result))
	}
}
