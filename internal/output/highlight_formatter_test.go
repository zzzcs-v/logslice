package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/logslice/internal/query"
)

func TestWriteHighlightedNoFilter(t *testing.T) {
	e := map[string]any{"level": "info", "msg": "started"}
	var buf bytes.Buffer
	if err := WriteHighlighted(&buf, e, nil, nil); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "level=info") {
		t.Errorf("expected level=info in output, got %q", out)
	}
}

func TestWriteHighlightedWithMatch(t *testing.T) {
	e := map[string]any{"level": "error", "msg": "disk full"}
	f := &query.Filter{Op: query.OpEq, Field: "level", Value: "error"}
	var buf bytes.Buffer
	if err := WriteHighlighted(&buf, e, f, nil); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "\033[1;33merror\033[0m") {
		t.Errorf("expected highlighted 'error', got %q", out)
	}
}

func TestWriteHighlightedSelectedFields(t *testing.T) {
	e := map[string]any{"level": "warn", "msg": "low memory", "ts": "2024-01-01"}
	var buf bytes.Buffer
	if err := WriteHighlighted(&buf, e, nil, []string{"msg"}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, "level=") {
		t.Errorf("expected level to be excluded, got %q", out)
	}
	if !strings.Contains(out, "msg=low memory") {
		t.Errorf("expected msg in output, got %q", out)
	}
}

func TestAnyToString(t *testing.T) {
	cases := []struct {
		input any
		want  string
	}{
		{"hello", "hello"},
		{float64(42), "42"},
		{true, "true"},
		{false, "false"},
		{nil, ""},
	}
	for _, c := range cases {
		got := anyToString(c.input)
		if got != c.want {
			t.Errorf("anyToString(%v) = %q, want %q", c.input, got, c.want)
		}
	}
}
