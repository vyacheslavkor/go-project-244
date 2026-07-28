package formatters

import (
	"code/internal/diff"
	"fmt"
	"strings"
)

type PlainFormatter struct{}

func (f *PlainFormatter) Format(d *diff.Diff) (string, error) {
	return processDiffPlain(d, ""), nil
}

func processDiffPlain(d *diff.Diff, path string) string {
	keys := d.ExtractKeys()

	lines := make([]string, 0, len(keys))

	for _, k := range keys {
		node := d.Nodes[k]

		switch node.Status {
		case "added":
			lines = append(lines, fmt.Sprintf("Property '%s%s' was added with value: %s", path, k, formatValueToPlain(node.Value)))
		case "removed":
			lines = append(lines, fmt.Sprintf("Property '%s%s' was removed", path, k))
		case "updated":
			lines = append(lines, fmt.Sprintf("Property '%s%s' was updated. From %s to %s", path, k, formatValueToPlain(node.OldValue), formatValueToPlain(node.Value)))
		case "unchanged":
			continue
		case "nested":
			lines = append(lines, processDiffPlain(node.Children, fmt.Sprintf("%s%s.", path, k)))
		}
	}

	return strings.Join(lines, "\n")
}

func formatValueToPlain(val any) string {
	if val == nil {
		return "null"
	}

	switch v := val.(type) {
	case string:
		return fmt.Sprintf("'%v'", v)
	case bool, int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprintf("%v", v)
	case map[string]any, []any:
		return "[complex value]"
	default:
		return fmt.Sprintf("%v", v)
	}
}
