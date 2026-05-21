package filter

import (
	"testing"

	"github.com/yourorg/logslice/internal/query"
)

func line(s string) []byte { return []byte(s) }

func TestMatchNilFilter(t *testing.T) {
	ok, err := Match(line(`{"level":"info"}`), nil)
	if err != nil || !ok {
		t.Fatal("nil filter should match everything")
	}
}

func TestMatchEq(t *testing.T) {
	f := &query.Filter{Field: "level", Op: query.OpEq, Value: "error"}
	ok, _ := Match(line(`{"level":"error"}`), f)
	if !ok {
		t.Fatal("expected match")
	}
	ok, _ = Match(line(`{"level":"info"}`), f)
	if ok {
		t.Fatal("expected no match")
	}
}

func TestMatchNeq(t *testing.T) {
	f := &query.Filter{Field: "level", Op: query.OpNeq, Value: "debug"}
	ok, _ := Match(line(`{"level":"info"}`), f)
	if !ok {
		t.Fatal("expected match")
	}
}

func TestMatchContains(t *testing.T) {
	f := &query.Filter{Field: "msg", Op: query.OpContains, Value: "timeout"}
	ok, _ := Match(line(`{"msg":"connection timeout occurred"}`), f)
	if !ok {
		t.Fatal("expected match")
	}
}

func TestMatchExists(t *testing.T) {
	f := &query.Filter{Field: "trace_id", Op: query.OpExists}
	ok, _ := Match(line(`{"trace_id":"abc123"}`), f)
	if !ok {
		t.Fatal("expected match")
	}
	ok, _ = Match(line(`{"level":"info"}`), f)
	if ok {
		t.Fatal("expected no match")
	}
}

func TestMatchGt(t *testing.T) {
	f := &query.Filter{Field: "status", Op: query.OpGt, Value: "400"}
	ok, _ := Match(line(`{"status":500}`), f)
	if !ok {
		t.Fatal("expected match")
	}
	ok, _ = Match(line(`{"status":200}`), f)
	if ok {
		t.Fatal("expected no match")
	}
}

func TestMatchInvalidJSON(t *testing.T) {
	f := &query.Filter{Field: "level", Op: query.OpEq, Value: "info"}
	_, err := Match(line(`not json`), f)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
