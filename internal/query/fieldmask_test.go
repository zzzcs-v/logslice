package query

import (
	"encoding/json"
	"testing"
)

func rawMaskEntry(t *testing.T, s string) json.RawMessage {
	t.Helper()
	return json.RawMessage(s)
}

func unmarshalMap(t *testing.T, raw json.RawMessage) map[string]json.RawMessage {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	return m
}

func TestFieldMaskNil(t *testing.T) {
	entry := rawMaskEntry(t, `{"level":"info","msg":"hello"}`)
	var fm *FieldMask
	out, err := fm.Apply(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(entry) {
		t.Errorf("expected unchanged entry, got %s", out)
	}
}

func TestFieldMaskInclude(t *testing.T) {
	entry := rawMaskEntry(t, `{"level":"info","msg":"hello","ts":"2024-01-01"}`)
	fm := &FieldMask{Include: []string{"level", "msg"}}
	out, err := fm.Apply(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := unmarshalMap(t, out)
	if _, ok := m["level"]; !ok {
		t.Error("expected 'level' field")
	}
	if _, ok := m["msg"]; !ok {
		t.Error("expected 'msg' field")
	}
	if _, ok := m["ts"]; ok {
		t.Error("expected 'ts' to be excluded")
	}
}

func TestFieldMaskExclude(t *testing.T) {
	entry := rawMaskEntry(t, `{"level":"info","msg":"hello","ts":"2024-01-01"}`)
	fm := &FieldMask{Exclude: []string{"ts"}}
	out, err := fm.Apply(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := unmarshalMap(t, out)
	if _, ok := m["ts"]; ok {
		t.Error("expected 'ts' to be excluded")
	}
	if _, ok := m["level"]; !ok {
		t.Error("expected 'level' to remain")
	}
}

func TestFieldMaskIncludeTakesPrecedence(t *testing.T) {
	entry := rawMaskEntry(t, `{"level":"info","msg":"hello","ts":"2024-01-01"}`)
	fm := &FieldMask{Include: []string{"level"}, Exclude: []string{"level"}}
	out, err := fm.Apply(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := unmarshalMap(t, out)
	if _, ok := m["level"]; !ok {
		t.Error("expected 'level' to be present when in Include")
	}
}

func TestFieldMaskInvalidJSON(t *testing.T) {
	entry := rawMaskEntry(t, `not-json`)
	fm := &FieldMask{Include: []string{"level"}}
	_, err := fm.Apply(entry)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
