package query

import (
	"testing"
)

func TestOptimizeNil(t *testing.T) {
	if Optimize(nil) != nil {
		t.Fatal("expected nil")
	}
}

func TestOptimizeLeaf(t *testing.T) {
	f := &Filter{Op: OpEq, Field: "level", Value: "info"}
	out := Optimize(f)
	if out.Op != OpEq || out.Field != "level" || out.Value != "info" {
		t.Fatalf("unexpected result: %+v", out)
	}
}

func TestOptimizeAndNilLeft(t *testing.T) {
	right := &Filter{Op: OpEq, Field: "level", Value: "info"}
	f := &Filter{Op: OpAnd, Left: nil, Right: right}
	out := Optimize(f)
	if out.Op != OpEq {
		t.Fatalf("expected leaf, got: %+v", out)
	}
}

func TestOptimizeAndNilRight(t *testing.T) {
	left := &Filter{Op: OpEq, Field: "level", Value: "info"}
	f := &Filter{Op: OpAnd, Left: left, Right: nil}
	out := Optimize(f)
	if out.Op != OpEq {
		t.Fatalf("expected leaf, got: %+v", out)
	}
}

func TestOptimizeOrBothPresent(t *testing.T) {
	left := &Filter{Op: OpEq, Field: "level", Value: "info"}
	right := &Filter{Op: OpEq, Field: "level", Value: "warn"}
	f := &Filter{Op: OpOr, Left: left, Right: right}
	out := Optimize(f)
	if out.Op != OpOr {
		t.Fatalf("expected OR, got: %+v", out)
	}
	if out.Left.Value != "info" || out.Right.Value != "warn" {
		t.Fatalf("unexpected children: %+v", out)
	}
}

func TestOptimizeAndBothPresent(t *testing.T) {
	left := &Filter{Op: OpExists, Field: "msg"}
	right := &Filter{Op: OpEq, Field: "level", Value: "error"}
	f := &Filter{Op: OpAnd, Left: left, Right: right}
	out := Optimize(f)
	if out.Op != OpAnd {
		t.Fatalf("expected AND, got: %+v", out)
	}
}
