package query

import "testing"

func TestDefaultPageOptions(t *testing.T) {
	opts := DefaultPageOptions()
	if opts.Offset != 0 || opts.Limit != 0 {
		t.Fatalf("expected zero values, got %+v", opts)
	}
}

func TestPaginatorNoLimitNoOffset(t *testing.T) {
	p := NewPaginator(DefaultPageOptions())
	for i := 0; i < 1000; i++ {
		emit, done := p.Accept()
		if !emit || done {
			t.Fatalf("entry %d: expected emit=true done=false", i)
		}
	}
	if p.Taken() != 1000 {
		t.Fatalf("expected 1000 taken, got %d", p.Taken())
	}
}

func TestPaginatorOffset(t *testing.T) {
	p := NewPaginator(PageOptions{Offset: 3, Limit: 0})
	for i := 0; i < 3; i++ {
		emit, done := p.Accept()
		if emit || done {
			t.Fatalf("entry %d should be skipped", i)
		}
	}
	emit, done := p.Accept()
	if !emit || done {
		t.Fatal("4th entry should be emitted")
	}
	if p.Taken() != 1 {
		t.Fatalf("expected taken=1, got %d", p.Taken())
	}
}

func TestPaginatorLimit(t *testing.T) {
	p := NewPaginator(PageOptions{Offset: 0, Limit: 2})

	emit, done := p.Accept()
	if !emit || done {
		t.Fatal("first entry: expected emit=true done=false")
	}
	emit, done = p.Accept()
	if !emit || done {
		t.Fatal("second entry: expected emit=true done=false")
	}
	// limit reached
	emit, done = p.Accept()
	if emit || !done {
		t.Fatal("third entry: expected emit=false done=true")
	}
	if p.Done() != true {
		t.Fatal("Done() should be true")
	}
}

func TestPaginatorOffsetAndLimit(t *testing.T) {
	p := NewPaginator(PageOptions{Offset: 2, Limit: 3})
	results := make([]bool, 7)
	for i := range results {
		emit, _ := p.Accept()
		results[i] = emit
	}
	// entries 0,1 skipped; 2,3,4 emitted; 5,6 rejected (done)
	expected := []bool{false, false, true, true, true, false, false}
	for i, want := range expected {
		if results[i] != want {
			t.Errorf("entry %d: emit=%v want=%v", i, results[i], want)
		}
	}
	if p.Taken() != 3 {
		t.Fatalf("expected taken=3, got %d", p.Taken())
	}
}

// TestPaginatorDoneStopsEmitting verifies that once the paginator is done,
// all subsequent Accept calls return emit=false and done=true.
func TestPaginatorDoneStopsEmitting(t *testing.T) {
	p := NewPaginator(PageOptions{Offset: 0, Limit: 1})

	emit, done := p.Accept()
	if !emit || done {
		t.Fatal("first entry: expected emit=true done=false")
	}

	for i := 0; i < 5; i++ {
		emit, done = p.Accept()
		if emit || !done {
			t.Fatalf("call %d after limit: expected emit=false done=true", i+2)
		}
	}
}
