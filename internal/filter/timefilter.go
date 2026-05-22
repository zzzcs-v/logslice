package filter

import (
	"encoding/json"
	"time"

	"github.com/user/logslice/internal/query"
)

// timeFields lists common timestamp keys found in structured logs.
var timeFields = []string{"time", "timestamp", "ts", "@timestamp"}

// MatchTimeRange reports whether a raw JSON log entry falls within tr.
// It checks a set of well-known timestamp fields and returns true when
// the first recognised field satisfies the range, or true when no
// timestamp field is found (pass-through).
func MatchTimeRange(raw json.RawMessage, tr query.TimeRange) bool {
	if tr.IsZero() {
		return true
	}

	var entry map[string]interface{}
	if err := json.Unmarshal(raw, &entry); err != nil {
		return true // malformed — let stream layer handle it
	}

	for _, field := range timeFields {
		v, ok := entry[field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		t, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			t, err = time.Parse(time.RFC3339, s)
			if err != nil {
				continue
			}
		}
		return tr.Contains(t)
	}

	// No recognised timestamp field — do not exclude.
	return true
}
