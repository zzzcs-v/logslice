package query

import (
	"encoding/json"
	"fmt"
	"sort"
)

// AggregateResult holds the count of log entries grouped by a field value.
type AggregateResult struct {
	Field  string
	Counts map[string]int
}

// Aggregate counts log entries grouped by the given field name.
// Each entry in entries is a raw JSON message.
func Aggregate(entries []json.RawMessage, field string) (*AggregateResult, error) {
	if field == "" {
		return nil, fmt.Errorf("aggregate: field name must not be empty")
	}

	counts := make(map[string]int)

	for _, raw := range entries {
		var obj map[string]interface{}
		if err := json.Unmarshal(raw, &obj); err != nil {
			// skip malformed entries
			continue
		}

		val, ok := obj[field]
		if !ok {
			counts["<missing>"]++
			continue
		}

		key := fmt.Sprintf("%v", val)
		counts[key]++
	}

	return &AggregateResult{Field: field, Counts: counts}, nil
}

// TopN returns up to n field values sorted by count descending.
type RankedEntry struct {
	Value string
	Count int
}

func (r *AggregateResult) TopN(n int) []RankedEntry {
	ranked := make([]RankedEntry, 0, len(r.Counts))
	for k, v := range r.Counts {
		ranked = append(ranked, RankedEntry{Value: k, Count: v})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Count != ranked[j].Count {
			return ranked[i].Count > ranked[j].Count
		}
		return ranked[i].Value < ranked[j].Value
	})
	if n > 0 && n < len(ranked) {
		return ranked[:n]
	}
	return ranked
}
