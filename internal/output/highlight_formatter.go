package output

import (
	"encoding/json"
	"io"

	"github.com/user/logslice/internal/query"
)

// WriteHighlighted writes a single log entry as pretty-printed text with
// matching terms highlighted using ANSI colour codes.
// fields controls which keys to emit; nil means all keys.
func WriteHighlighted(w io.Writer, entry map[string]any, f *query.Filter, fields []string) error {
	keys := fields
	if len(keys) == 0 {
		keys = sortedKeys(entry)
	}

	for _, k := range keys {
		v, ok := entry[k]
		if !ok {
			continue
		}
		raw := anyToString(v)
		highlighted := query.HighlightMatch(f, raw)
		if _, err := io.WriteString(w, k+"="+highlighted+" "); err != nil {
			return err
		}
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// anyToString converts an arbitrary JSON value to a plain string for display.
func anyToString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		b, _ := json.Marshal(t)
		return string(b)
	case bool:
		if t {
			return "true"
		}
		return "false"
	case nil:
		return ""
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
