package parsers

import "time"

// normalizeMap converts parser-specific values in a nested map
// to their common representation for comparison.
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
		return normalizeTime(v)
	default:
		return v
	}
}

func normalizeTime(t time.Time) string {
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		return t.Format("2006-01-02")
	}

	return t.Format(time.RFC3339)
}
