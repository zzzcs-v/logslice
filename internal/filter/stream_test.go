package filter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/query"
)

const sampleLogs = `{"level":"info","msg":"server started"}
{"level":"error","msg":"connection refused"}
{"level":"error","msg":"timeout"}
{"level":"debug","msg":"heartbeat"}
`

func TestStreamFilterMatchAll(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	var w bytes.Buffer
	matched, total, err := StreamFilter(r, &w, nil)
	if err != nil {
		t.Fatal(err)
	}
	if total != 4 || matched != 4 {
		t.Fatalf("expected 4/4, got %d/%d", matched, total)
	}
}

func TestStreamFilterByLevel(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	var w bytes.Buffer
	f := &query.Filter{Field: "level", Op: query.OpEq, Value: "error"}
	matched, total, err := StreamFilter(r, &w, f)
	if err != nil {
		t.Fatal(err)
	}
	if total != 4 {
		t.Fatalf("expected total=4, got %d", total)
	}
	if matched != 2 {
		t.Fatalf("expected matched=2, got %d", matched)
	}
	lines := strings.Split(strings.TrimSpace(w.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(lines))
	}
}

func TestStreamFilterSkipMalformed(t *testing.T) {
	input := "{\"level\":\"info\"}\nnot-json\n{\"level\":\"info\"}\n"
	r := strings.NewReader(input)
	var w bytes.Buffer
	f := &query.Filter{Field: "level", Op: query.OpEq, Value: "info"}
	matched, total, err := StreamFilter(r, &w, f)
	if err != nil {
		t.Fatal(err)
	}
	if total != 3 || matched != 2 {
		t.Fatalf("expected 2/3, got %d/%d", matched, total)
	}
}
