package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format controls how matched log lines are rendered.
type Format string

const (
	FormatJSON   Format = "json"
	FormatPretty Format = "pretty"
	FormatText   Format = "text"
)

// Formatter writes log entries to w in the chosen format.
type Formatter struct {
	Writer io.Writer
	Format Format
	// Fields limits output to these keys when using FormatText; empty means all.
	Fields []string
}

// Write renders a single parsed log entry.
func (f *Formatter) Write(entry map[string]any) error {
	switch f.Format {
	case FormatPretty:
		return f.writePretty(entry)
	case FormatText:
		return f.writeText(entry)
	default: // FormatJSON
		return f.writeJSON(entry)
	}
}

func (f *Formatter) writeJSON(entry map[string]any) error {
	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f.Writer, string(b))
	return err
}

func (f *Formatter) writePretty(entry map[string]any) error {
	b, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f.Writer, string(b))
	return err
}

func (f *Formatter) writeText(entry map[string]any) error {
	keys := f.Fields
	if len(keys) == 0 {
		keys = sortedKeys(entry)
	}
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		v, ok := entry[k]
		if !ok {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	_, err := fmt.Fprintln(f.Writer, strings.Join(parts, " "))
	return err
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// stable order: level, time, msg first, then the rest alphabetically
	priority := map[string]int{"level": 0, "time": 1, "msg": 2}
	sort(keys, priority)
	return keys
}

func sort(keys []string, priority map[string]int) {
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0; j-- {
			pi, oki := priority[keys[j]]
			pj, okj := priority[keys[j-1]]
			swap := false
			if oki && okj {
				swap = pi < pj
			} else if oki {
				swap = true
			} else if !okj {
				swap = keys[j] < keys[j-1]
			}
			if swap {
				keys[j], keys[j-1] = keys[j-1], keys[j]
			} else {
				break
			}
		}
	}
}
