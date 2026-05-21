package query

import (
	"testing"
)

func TestFilterStringNil(t *testing.T) {
	var f *Filter
	if got := f.String(); got != "<nil>" {
		t.Errorf("expected <nil>, got %q", got)
	}
}

func TestFilterStringEq(t *testing.T) {
	f := &Filter{Field: "level", Op: OpEq, Value: "error"}
	expected := "level = error"
	if got := f.String(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFilterStringExists(t *testing.T) {
	f := &Filter{Field: "trace_id", Op: OpExists}
	expected := "trace_id exists"
	if got := f.String(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFilterStringWithAnd(t *testing.T) {
	f := &Filter{
		Field: "level",
		Op:    OpEq,
		Value: "error",
		And:   &Filter{Field: "service", Op: OpEq, Value: "api"},
	}
	expected := "level = error AND service = api"
	if got := f.String(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFilterStringWithOr(t *testing.T) {
	f := &Filter{
		Field: "level",
		Op:    OpEq,
		Value: "error",
		Or:    &Filter{Field: "level", Op: OpEq, Value: "warn"},
	}
	expected := "level = error OR level = warn"
	if got := f.String(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestOperatorIsComparison(t *testing.T) {
	cases := []struct {
		op      Operator
		want    bool
	}{
		{OpEq, true},
		{OpNeq, true},
		{OpContains, true},
		{OpGt, true},
		{OpLt, true},
		{OpGte, true},
		{OpLte, true},
		{OpExists, false},
	}
	for _, c := range cases {
		if got := c.op.IsComparison(); got != c.want {
			t.Errorf("op %q: expected IsComparison=%v, got %v", c.op, c.want, got)
		}
	}
}
