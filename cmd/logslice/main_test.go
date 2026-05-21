package main

import (
	"testing"
)

func TestSplitFieldsEmpty(t *testing.T) {
	result := splitFields("")
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %v", result)
	}
}

func TestSplitFieldsSingle(t *testing.T) {
	result := splitFields("level")
	if len(result) != 1 || result[0] != "level" {
		t.Fatalf("expected [level], got %v", result)
	}
}

func TestSplitFieldsMultiple(t *testing.T) {
	result := splitFields("level,msg,time")
	expected := []string{"level", "msg", "time"}
	if len(result) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("index %d: expected %q, got %q", i, v, result[i])
		}
	}
}

func TestSplitFieldsTrailingComma(t *testing.T) {
	result := splitFields("level,msg,")
	if len(result) != 2 {
		t.Fatalf("expected 2 fields, got %v", result)
	}
}

func TestSplitFieldsLeadingComma(t *testing.T) {
	result := splitFields(",level")
	if len(result) != 1 || result[0] != "level" {
		t.Fatalf("expected [level], got %v", result)
	}
}
