package query

import (
	"encoding/json"
	"sort"
)

// GroupResult holds the entries for a single group key.
type GroupResult struct {
	Key     string
	Entries []json.RawMessage
}

// GroupBy partitions entries by the distinct values of the given field.
// Entries missing the field are collected under the key "".
// Results are returned sorted by key ascending.
func GroupBy(entries []json.RawMessage, field string) []GroupResult {
	if field == "" || len(entries) == 0 {
		return []GroupResult{{Key: "", Entries: entries}}
	}

	index := make(map[string][]json.RawMessage)
	order := []string{}

	for _, raw := range entries {
		var obj map[string]interface{}
		if err := json.Unmarshal(raw, &obj); err != nil {
			if _, exists := index[""]; !exists {
				order = append(order, "")
			}
			index[""] = append(index[""], raw)
			continue
		}

		key := ""
		if v, ok := obj[field]; ok {
			key = fieldToString(v)
		}

		if _, exists := index[key]; !exists {
			order = append(order, key)
		}
		index[key] = append(index[key], raw)
	}

	sort.Strings(order)

	results := make([]GroupResult, 0, len(order))
	for _, k := range order {
		results = append(results, GroupResult{Key: k, Entries: index[k]})
	}
	return results
}
