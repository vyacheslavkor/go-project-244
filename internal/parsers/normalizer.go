package parsers

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
	default:
		return v
	}
}
