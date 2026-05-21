package query

import (
	"testing"
)

func TestValidateNil(t *testing.T) {
	if err := Validate(nil); err != nil {
		t.Fatalf("expected nil error for nil filter, got %v", err)
	}
}

func TestValidateEq(t *testing.T) {
	f := &Filter{Op: OpEq, Field: "level", Value: "info"}
	if err := Validate(f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateEqEmptyField(t *testing.T) {
	f := &Filter{Op: OpEq, Field: "", Value: "info"}
	if err := Validate(f); err == nil {
		t.Fatal("expected error for empty field name, got nil")
	}
}

func TestValidateExists(t *testing.T) {
	f := &Filter{Op: OpExists, Field: "request_id"}
	if err := Validate(f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateExistsEmptyField(t *testing.T) {
	f := &Filter{Op: OpExists, Field: "   "}
	if err := Validate(f); err == nil {
		t.Fatal("expected error for blank field name, got nil")
	}
}

func TestValidateAndMissingRight(t *testing.T) {
	f := &Filter{Op: OpAnd, Left: &Filter{Op: OpExists, Field: "x"}, Right: nil}
	if err := Validate(f); err == nil {
		t.Fatal("expected error when right operand is missing")
	}
}

func TestValidateNotMissingOperand(t *testing.T) {
	f := &Filter{Op: OpNot, Left: nil}
	if err := Validate(f); err == nil {
		t.Fatal("expected error when NOT has no operand")
	}
}

func TestValidateNestedAnd(t *testing.T) {
	f := &Filter{
		Op: OpAnd,
		Left:  &Filter{Op: OpEq, Field: "level", Value: "error"},
		Right: &Filter{Op: OpContains, Field: "msg", Value: "timeout"},
	}
	if err := Validate(f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateUnknownOp(t *testing.T) {
	f := &Filter{Op: "BETWEEN", Field: "age"}
	if err := Validate(f); err == nil {
		t.Fatal("expected error for unknown operator")
	}
}
