package formatters

import (
	"code/internal/diff"
	"fmt"
	"strings"
)

type PlainFormatter struct{}

func (f *PlainFormatter) Format(d *diff.Diff) string {
	return processDiffPlain(d, "")
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

	_, isBool := val.(bool)
	if isBool {
		return fmt.Sprintf("%v", val)
	}

	_, isMap := val.(map[string]any)
	if isMap {
		return "[complex value]"
	}

	return fmt.Sprintf("'%v'", val)
}
