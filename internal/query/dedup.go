package query

import (
	"encoding/json"
	"fmt"
)

// DedupEntries removes duplicate log entries based on the value of a specific field.
// The first occurrence of each unique field value is kept; subsequent duplicates are dropped.
// If field is empty, all entries are returned as-is.
func DedupEntries(entries []json.RawMessage, field string) ([]json.RawMessage, error) {
	if field == "" {
		return entries, nil
	}

	seen := make(map[string]struct{})
	result := make([]json.RawMessage, 0, len(entries))

	for _, raw := range entries {
		var m map[string]interface{}
		if err := json.Unmarshal(raw, &m); err != nil {
			// Skip malformed entries
			continue
		}

		key, err := fieldKey(m, field)
		if err != nil {
			// Field missing — treat as unique key per entry index to preserve it
			result = append(result, raw)
			continue
		}

		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		result = append(result, raw)
	}

	return result, nil
}

// fieldKey extracts a string key from a map for a given field name.
func fieldKey(m map[string]interface{}, field string) (string, error) {
	val, ok := m[field]
	if !ok {
		return "", fmt.Errorf("field %q not found", field)
	}
	switch v := val.(type) {
	case string:
		return v, nil
	case float64:
		return fmt.Sprintf("%g", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	case nil:
		return "<nil>", nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}
