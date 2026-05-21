package query

// Operator represents a comparison or existence operator in a filter expression.
type Operator string

const (
	OpEq       Operator = "="
	OpNeq      Operator = "!="
	OpContains Operator = "~"
	OpExists   Operator = "exists"
	OpGt       Operator = ">"
	OpLt       Operator = "<"
	OpGte      Operator = ">="
	OpLte      Operator = "<="
)

// Filter represents a single filter expression: field op value.
type Filter struct {
	Field    string
	Op       Operator
	Value    string
	And      *Filter
	Or       *Filter
}

// Query is the top-level parsed query structure.
type Query struct {
	Filter *Filter
}

// String returns a human-readable representation of the filter for debugging.
func (f *Filter) String() string {
	if f == nil {
		return "<nil>"
	}
	s := string(f.Field) + " " + string(f.Op)
	if f.Op != OpExists {
		s += " " + f.Value
	}
	if f.And != nil {
		s += " AND " + f.And.String()
	}
	if f.Or != nil {
		s += " OR " + f.Or.String()
	}
	return s
}

// IsComparison returns true if the operator compares a value (not just checks existence).
func (op Operator) IsComparison() bool {
	switch op {
	case OpEq, OpNeq, OpContains, OpGt, OpLt, OpGte, OpLte:
		return true
	}
	return false
}
