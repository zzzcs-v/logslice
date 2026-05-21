package filter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/yourorg/logslice/internal/query"
)

// Match checks whether a JSON log line satisfies the given filter.
func Match(line []byte, f *query.Filter) (bool, error) {
	if f == nil {
		return true, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal(line, &record); err != nil {
		return false, err
	}

	return matchFilter(record, f), nil
}

func matchFilter(record map[string]interface{}, f *query.Filter) bool {
	val, exists := record[f.Field]

	switch f.Op {
	case query.OpExists:
		return exists
	case query.OpEq:
		return exists && toString(val) == f.Value
	case query.OpNeq:
		return !exists || toString(val) != f.Value
	case query.OpContains:
		return exists && strings.Contains(toString(val), f.Value)
	case query.OpGt:
		return exists && compareNumeric(val, f.Value) > 0
	case query.OpLt:
		return exists && compareNumeric(val, f.Value) < 0
	}
	return false
}

func toString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", t), "0"), ".")
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func compareNumeric(val interface{}, target string) int {
	targetF, err := strconv.ParseFloat(target, 64)
	if err != nil {
		return 0
	}
	var valF float64
	switch t := val.(type) {
	case float64:
		valF = t
	case string:
		valF, err = strconv.ParseFloat(t, 64)
		if err != nil {
			return 0
		}
	default:
		return 0
	}
	if valF > targetF {
		return 1
	} else if valF < targetF {
		return -1
	}
	return 0
}
