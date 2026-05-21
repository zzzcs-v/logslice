package query

import "strings"

// HighlightMatch returns a copy of the input string with any occurrences of
// the search terms from the filter wrapped in ANSI escape codes.
func HighlightMatch(f *Filter, value string) string {
	if f == nil || value == "" {
		return value
	}
	terms := collectTerms(f)
	result := value
	for _, term := range terms {
		if term == "" {
			continue
		}
		result = replaceCI(result, term, "\033[1;33m"+term+"\033[0m")
	}
	return result
}

// collectTerms walks the filter tree and gathers literal string values used in
// Eq, Neq, and Contains predicates so we know what to highlight.
func collectTerms(f *Filter) []string {
	if f == nil {
		return nil
	}
	var terms []string
	switch f.Op {
	case OpEq, OpNeq, OpContains:
		if f.Value != "" {
			terms = append(terms, f.Value)
		}
	case OpAnd, OpOr:
		terms = append(terms, collectTerms(f.Left)...)
		terms = append(terms, collectTerms(f.Right)...)
	}
	return terms
}

// replaceCI replaces all case-insensitive occurrences of old in s with new.
func replaceCI(s, old, replacement string) string {
	if old == "" {
		return s
	}
	lower := strings.ToLower(s)
	lowerOld := strings.ToLower(old)
	var result strings.Builder
	for {
		idx := strings.Index(lower, lowerOld)
		if idx < 0 {
			result.WriteString(s)
			break
		}
		result.WriteString(s[:idx])
		result.WriteString(replacement)
		s = s[idx+len(old):]
		lower = lower[idx+len(old):]
	}
	return result.String()
}
