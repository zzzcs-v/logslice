package query

import (
	"testing"
)

func TestParseEmpty(t *testing.T) {
	q, err := Parse("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(q.Filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(q.Filters))
	}
}

func TestParseEq(t *testing.T) {
	q, err := Parse("level=error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(q.Filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(q.Filters))
	}
	f := q.Filters[0]
	assertFilter(t, f, "level", OpEq, "error")
}

func TestParseNeq(t *testing.T) {
	q, err := Parse("level!=debug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertFilter(t, q.Filters[0], "level", OpNeq, "debug")
}

func TestParseContains(t *testing.T) {
	q, err := Parse("msg~timeout")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertFilter(t, q.Filters[0], "msg", OpContains, "timeout")
}

func TestParseExists(t *testing.T) {
	q, err := Parse("trace_id?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertFilter(t, q.Filters[0], "trace_id", OpExists, "")
}

func TestParseMultiple(t *testing.T) {
	q, err := Parse("level=error service=auth msg~timeout")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(q.Filters) != 3 {
		t.Fatalf("expected 3 filters, got %d", len(q.Filters))
	}
	assertFilter(t, q.Filters[0], "level", OpEq, "error")
	assertFilter(t, q.Filters[1], "service", OpEq, "auth")
	assertFilter(t, q.Filters[2], "msg", OpContains, "timeout")
}

func TestParseInvalidExpression(t *testing.T) {
	_, err := Parse("justaplainword")
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestParseEmptyField(t *testing.T) {
	_, err := Parse("=value")
	if err == nil {
		t.Fatal("expected error for empty field name, got nil")
	}
}

func assertFilter(t *testing.T, f Filter, field string, op Op, value string) {
	t.Helper()
	if f.Field != field {
		t.Errorf("field: expected %q, got %q", field, f.Field)
	}
	if f.Op != op {
		t.Errorf("op: expected %q, got %q", op, f.Op)
	}
	if f.Value != value {
		t.Errorf("value: expected %q, got %q", value, f.Value)
	}
}
