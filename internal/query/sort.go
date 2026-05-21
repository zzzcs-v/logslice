package query

import (
	"encoding/json"
	"sort"
	"strconv"
)

// SortOrder represents ascending or descending sort direction.
type SortOrder int

const (
	Ascending  SortOrder = iota
	Descending SortOrder = iota
)

// SortOptions configures sorting of log entries.
type SortOptions struct {
	Field string
	Order SortOrder
}

// SortEntries sorts a slice of raw JSON log entries by the given field.
// Entries missing the field are placed at the end.
func SortEntries(entries []json.RawMessage, opts SortOptions) []json.RawMessage {
	if opts.Field == "" || len(entries) == 0 {
		return entries
	}

	type indexed struct {
		raw   json.RawMessage
		value string
		missing bool
	}

	indexed := make([]struct {
		raw     json.RawMessage
		value   string
		missing bool
	}, len(entries))

	for i, raw := range entries {
		var obj map[string]interface{}
		indexed[i].raw = raw
		if err := json.Unmarshal(raw, &obj); err != nil {
			indexed[i].missing = true
			continue
		}
		v, ok := obj[opts.Field]
		if !ok {
			indexed[i].missing = true
			continue
		}
		indexed[i].value = fieldToString(v)
	}

	sort.SliceStable(indexed, func(i, j int) bool {
		if indexed[i].missing && indexed[j].missing {
			return false
		}
		if indexed[i].missing {
			return false
		}
		if indexed[j].missing {
			return true
		}
		less := indexed[i].value < indexed[j].value
		if opts.Order == Descending {
			return !less
		}
		return less
	})

	result := make([]json.RawMessage, len(entries))
	for i, item := range indexed {
		result[i] = item.raw
	}
	return result
}

func fieldToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}
