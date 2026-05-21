package output

import (
	"bytes"
	"strings"
	"testing"
)

func entry() map[string]any {
	return map[string]any{
		"level": "info",
		"msg":   "hello world",
		"time":  "2024-01-01T00:00:00Z",
		"svc":   "api",
	}
}

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatJSON}
	if err := f.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"level":"info"`) {
		t.Errorf("expected level field in JSON output, got: %s", out)
	}
	if !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected JSON object, got: %s", out)
	}
}

func TestWritePretty(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatPretty}
	if err := f.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\n") {
		t.Errorf("pretty format should contain newlines, got: %s", out)
	}
	if !strings.Contains(out, `"msg": "hello world"`) {
		t.Errorf("expected indented msg field, got: %s", out)
	}
}

func TestWriteTextAllFields(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatText}
	if err := f.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"level=info", "msg=hello world", "svc=api"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in text output, got: %s", want, out)
		}
	}
}

func TestWriteTextSelectedFields(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatText, Fields: []string{"level", "msg"}}
	if err := f.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if out != "level=info msg=hello world" {
		t.Errorf("unexpected text output: %q", out)
	}
}

func TestWriteTextMissingField(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatText, Fields: []string{"level", "missing"}}
	if err := f.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if out != "level=info" {
		t.Errorf("expected only present fields, got: %q", out)
	}
}

func TestSortedKeysPriority(t *testing.T) {
	m := map[string]any{"svc": "x", "msg": "m", "time": "t", "level": "l", "aaa": "a"}
	keys := sortedKeys(m)
	if keys[0] != "level" || keys[1] != "time" || keys[2] != "msg" {
		t.Errorf("expected level,time,msg first, got: %v", keys)
	}
}
