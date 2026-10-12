package parsers

import "time"

var timeLayouts = []string{
	"2006-1-2T15:4:5.999999999Z07:00",
	"2006-1-2 15:4:5.999999999",
	"2006-1-2",
}

// normalizeMap converts values in a nested map to the form they are
// compared in, so that equal values from JSON and YAML match.
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
	case string:
		return normalizeString(v)
	case time.Time:
		return normalizeTime(v)
	default:
		return v
	}
}

func normalizeString(s string) string {
	for _, layout := range timeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return normalizeTime(t)
		}
	}

	return s
}

func normalizeTime(t time.Time) string {
	if _, offset := t.Zone(); offset != 0 || t.Nanosecond() != 0 {
		return t.Format(time.RFC3339Nano)
	}

	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		return t.Format(time.DateOnly)
	}

	return t.Format(time.DateTime)
}
