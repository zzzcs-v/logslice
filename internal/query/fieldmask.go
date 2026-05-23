package query

import (
	"encoding/json"
	"strings"
)

// FieldMask selects or excludes specific fields from a log entry.
type FieldMask struct {
	Include []string
	Exclude []string
}

// Apply returns a new raw message containing only the desired fields.
// If Include is non-empty, only those fields are kept.
// If Exclude is non-empty, those fields are removed.
// Include takes precedence over Exclude if both are set.
func (fm *FieldMask) Apply(entry json.RawMessage) (json.RawMessage, error) {
	if fm == nil || (len(fm.Include) == 0 && len(fm.Exclude) == 0) {
		return entry, nil
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(entry, &m); err != nil {
		return entry, err
	}

	result := make(map[string]json.RawMessage)

	if len(fm.Include) > 0 {
		includeSet := toSet(fm.Include)
		for k, v := range m {
			if includeSet[strings.ToLower(k)] {
				result[k] = v
			}
		}
	} else {
		excludeSet := toSet(fm.Exclude)
		for k, v := range m {
			if !excludeSet[strings.ToLower(k)] {
				result[k] = v
			}
		}
	}

	out, err := json.Marshal(result)
	if err != nil {
		return entry, err
	}
	return out, nil
}

func toSet(fields []string) map[string]bool {
	s := make(map[string]bool, len(fields))
	for _, f := range fields {
		s[strings.ToLower(f)] = true
	}
	return s
}
