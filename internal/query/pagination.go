package query

// PageOptions controls how many log entries to skip and how many to return.
type PageOptions struct {
	Offset int // number of matching entries to skip
	Limit  int // max entries to return; 0 means unlimited
}

// DefaultPageOptions returns options that return all entries.
func DefaultPageOptions() PageOptions {
	return PageOptions{Offset: 0, Limit: 0}
}

// Paginator tracks state while streaming entries through a page window.
type Paginator struct {
	opts    PageOptions
	skipped int
	taken   int
}

// NewPaginator creates a Paginator from the given options.
func NewPaginator(opts PageOptions) *Paginator {
	return &Paginator{opts: opts}
}

// Accept decides whether the next matching entry should be emitted.
// It must be called once per entry that passed the filter.
//
// Returns:
//   - emit=true  → caller should output this entry
//   - done=true  → limit reached, caller should stop reading
func (p *Paginator) Accept() (emit bool, done bool) {
	if p.skipped < p.opts.Offset {
		p.skipped++
		return false, false
	}

	if p.opts.Limit > 0 && p.taken >= p.opts.Limit {
		return false, true
	}

	p.taken++
	return true, false
}

// Taken returns how many entries have been emitted so far.
func (p *Paginator) Taken() int { return p.taken }

// Done reports whether the limit has been reached.
func (p *Paginator) Done() bool {
	return p.opts.Limit > 0 && p.taken >= p.opts.Limit
}
