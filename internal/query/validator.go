package query

import (
	"fmt"
	"strings"
)

// ValidationError describes a problem found during query validation.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("query validation error: %s", e.Message)
}

// Validate checks a parsed Filter tree for semantic correctness.
// It returns a ValidationError if any issue is found, or nil on success.
func Validate(f *Filter) error {
	if f == nil {
		return nil
	}
	return validateFilter(f)
}

func validateFilter(f *Filter) error {
	switch f.Op {
	case OpAnd, OpOr:
		if f.Left == nil || f.Right == nil {
			return &ValidationError{Message: fmt.Sprintf("operator %q requires both left and right operands", f.Op)}
		}
		if err := validateFilter(f.Left); err != nil {
			return err
		}
		return validateFilter(f.Right)

	case OpNot:
		if f.Left == nil {
			return &ValidationError{Message: "NOT operator requires an operand"}
		}
		return validateFilter(f.Left)

	case OpEq, OpNeq, OpContains, OpGt, OpLt, OpGte, OpLte:
		if strings.TrimSpace(f.Field) == "" {
			return &ValidationError{Message: fmt.Sprintf("operator %q requires a non-empty field name", f.Op)}
		}
		return nil

	case OpExists:
		if strings.TrimSpace(f.Field) == "" {
			return &ValidationError{Message: "EXISTS requires a non-empty field name"}
		}
		return nil

	default:
		return &ValidationError{Message: fmt.Sprintf("unknown operator %q", f.Op)}
	}
}
