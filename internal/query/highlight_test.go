package query

import (
	"strings"
	"testing"
)

func TestHighlightNilFilter(t *testing.T) {
	out := HighlightMatch(nil, "hello world")
	if out != "hello world" {
		t.Fatalf("expected unchanged string, got %q", out)
	}
}

func TestHighlightEmptyValue(t *testing.T) {
	f := &Filter{Op: OpEq, Field: "level", Value: "error"}
	out := HighlightMatch(f, "")
	if out != "" {
		t.Fatalf("expected empty string, got %q", out)
	}
}

func TestHighlightEqMatch(t *testing.T) {
	f := &Filter{Op: OpEq, Field: "level", Value: "error"}
	out := HighlightMatch(f, "level is error here")
	if !strings.Contains(out, "\033[1;33merror\033[0m") {
		t.Fatalf("expected ANSI highlight, got %q", out)
	}
}

func TestHighlightCaseInsensitive(t *testing.T) {
	f := &Filter{Op: OpContains, Field: "msg", Value: "warn"}
	out := HighlightMatch(f, "WARNING: disk full")
	// original casing is preserved in the surrounding text
	if !strings.Contains(out, "\033[1;33m") {
		t.Fatalf("expected highlight for case-insensitive match, got %q", out)
	}
}

func TestHighlightAndFilter(t *testing.T) {
	f := &Filter{
		Op: OpAnd,
		Left:  &Filter{Op: OpEq, Field: "level", Value: "error"},
		Right: &Filter{Op: OpContains, Field: "msg", Value: "timeout"},
	}
	out := HighlightMatch(f, "connection timeout error")
	if !strings.Contains(out, "\033[1;33merror\033[0m") {
		t.Fatalf("expected 'error' highlighted, got %q", out)
	}
	if !strings.Contains(out, "\033[1;33mtimeout\033[0m") {
		t.Fatalf("expected 'timeout' highlighted, got %q", out)
	}
}

func TestHighlightExistsNoTerm(t *testing.T) {
	f := &Filter{Op: OpExists, Field: "trace_id"}
	out := HighlightMatch(f, "some log line")
	// Exists has no value term, string should be unchanged
	if out != "some log line" {
		t.Fatalf("expected unchanged string for Exists op, got %q", out)
	}
}

func TestCollectTerms(t *testing.T) {
	f := &Filter{
		Op: OpOr,
		Left:  &Filter{Op: OpEq, Field: "a", Value: "foo"},
		Right: &Filter{Op: OpNeq, Field: "b", Value: "bar"},
	}
	terms := collectTerms(f)
	if len(terms) != 2 {
		t.Fatalf("expected 2 terms, got %d: %v", len(terms), terms)
	}
}
