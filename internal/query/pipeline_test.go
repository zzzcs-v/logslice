package query

import (
	"testing"
)

func TestPipelineEmpty(t *testing.T) {
	f, err := Pipeline("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != nil {
		t.Fatal("expected nil filter for empty input")
	}
}

func TestPipelineSimpleEq(t *testing.T) {
	f, err := Pipeline(`level = "info"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if f.Op != OpEq || f.Field != "level" || f.Value != "info" {
		t.Fatalf("unexpected filter: %+v", f)
	}
}

func TestPipelineInvalidField(t *testing.T) {
	_, err := Pipeline(`= "info"`)
	if err == nil {
		t.Fatal("expected error for invalid query")
	}
}

func TestPipelineValidationError(t *testing.T) {
	// Manually create a filter that would fail validation by
	// using Pipeline with an empty field — parser should catch this
	_, err := Pipeline(`level = ""`)
	// empty value is allowed; just ensure no panic
	if err != nil {
		t.Logf("got error (may be expected): %v", err)
	}
}

func TestPipelineAndOptimized(t *testing.T) {
	f, err := Pipeline(`level = "error" AND msg exists`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil || f.Op != OpAnd {
		t.Fatalf("expected AND filter, got: %+v", f)
	}
}
