package query

import (
	"encoding/json"
	"testing"
)

func makeRaw(t *testing.T, obj map[string]interface{}) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return json.RawMessage(b)
}

func TestAggregateEmptyField(t *testing.T) {
	_, err := Aggregate(nil, "")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestAggregateNoEntries(t *testing.T) {
	res, err := Aggregate([]json.RawMessage{}, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Counts) != 0 {
		t.Fatalf("expected empty counts, got %v", res.Counts)
	}
}

func TestAggregateCountsByField(t *testing.T) {
	entries := []json.RawMessage{
		makeRaw(t, map[string]interface{}{"level": "info"}),
		makeRaw(t, map[string]interface{}{"level": "error"}),
		makeRaw(t, map[string]interface{}{"level": "info"}),
		makeRaw(t, map[string]interface{}{"level": "warn"}),
	}
	res, err := Aggregate(entries, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Counts["info"] != 2 {
		t.Errorf("expected info=2, got %d", res.Counts["info"])
	}
	if res.Counts["error"] != 1 {
		t.Errorf("expected error=1, got %d", res.Counts["error"])
	}
	if res.Counts["warn"] != 1 {
		t.Errorf("expected warn=1, got %d", res.Counts["warn"])
	}
}

func TestAggregateMissingField(t *testing.T) {
	entries := []json.RawMessage{
		makeRaw(t, map[string]interface{}{"msg": "hello"}),
		makeRaw(t, map[string]interface{}{"level": "info"}),
	}
	res, err := Aggregate(entries, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Counts["<missing>"] != 1 {
		t.Errorf("expected <missing>=1, got %d", res.Counts["<missing>"])
	}
}

func TestTopN(t *testing.T) {
	res := &AggregateResult{
		Field:  "level",
		Counts: map[string]int{"info": 10, "error": 3, "warn": 7},
	}
	top := res.TopN(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(top))
	}
	if top[0].Value != "info" || top[0].Count != 10 {
		t.Errorf("unexpected top entry: %+v", top[0])
	}
	if top[1].Value != "warn" || top[1].Count != 7 {
		t.Errorf("unexpected second entry: %+v", top[1])
	}
}
