package query

import (
	"fmt"
	"strings"
)

// Op represents a comparison operator in a filter expression.
type Op string

const (
	OpEq      Op = "="
	OpNeq     Op = "!="
	OpContains Op = "~"
	OpExists   Op = "?"
)

// Filter represents a single parsed filter expression like `level=error` or `msg~timeout`.
type Filter struct {
	Field string
	Op    Op
	Value string
}

// Query holds the full parsed query with one or more filters.
type Query struct {
	Filters []Filter
}

// Parse parses a query string into a Query.
// Expressions are space-separated. Supported forms:
//   field=value   — exact match
//   field!=value  — not equal
//   field~value   — contains (substring)
//   field?        — field exists
func Parse(raw string) (*Query, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return &Query{}, nil
	}

	parts := strings.Fields(raw)
	filters := make([]Filter, 0, len(parts))

	for _, part := range parts {
		f, err := parseFilter(part)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}

	return &Query{Filters: filters}, nil
}

func parseFilter(expr string) (Filter, error) {
	// Order matters: check != before =
	for _, op := range []Op{OpNeq, OpEq, OpContains} {
		if idx := strings.Index(expr, string(op)); idx > 0 {
			field := expr[:idx]
			value := expr[idx+len(op):]
			if field == "" {
				return Filter{}, fmt.Errorf("empty field name in expression %q", expr)
			}
			return Filter{Field: field, Op: op, Value: value}, nil
		}
	}

	// Check exists operator: field?
	if strings.HasSuffix(expr, string(OpExists)) {
		field := strings.TrimSuffix(expr, string(OpExists))
		if field == "" {
			return Filter{}, fmt.Errorf("empty field name in expression %q", expr)
		}
		return Filter{Field: field, Op: OpExists}, nil
	}

	return Filter{}, fmt.Errorf("invalid filter expression %q: no operator found", expr)
}
