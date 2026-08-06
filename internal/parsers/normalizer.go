package parsers

import "time"

// normalizeMap normalizes numeric values in a nested map to float64,
// matching encoding/json so YAML ints and floats compare equally in diffs.
func normalizeMap(m map[string]any) {
	for k, v := range m {
		m[k] = normalizeValue(v)
	}
}

func normalizeValue(val any) any {
	switch v := val.(type) {
	case int:
		return float64(v)
	case uint64:
		return float64(v)
	case map[string]any:
		normalizeMap(v)
		return v
	case []any:
		for i, item := range v {
			v[i] = normalizeValue(item)
		}
		return v
	case time.Time:
		if v.Hour() == 0 && v.Minute() == 0 && v.Second() == 0 {
			return v.Format("2006-01-02")
		}
		return v.Format(time.RFC3339)
	default:
		return v
	}
}
