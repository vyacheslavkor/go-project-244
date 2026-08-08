package formatters

import (
	"code/internal/diff"
	"fmt"
	"slices"
	"strings"
)

// Stylish indent model (keys share one column at each depth):
//
//	indentStep   — spaces added per nesting level for the key column
//	markerWidth  — width of "+ " / "- " (or padding spaces when unmarked)
//
// At depth d, marked lines get (indentStep*d - markerWidth) leading spaces
// before the marker; unmarked lines get indentStep*d spaces before the key.
const (
	indentStep  = 4
	markerWidth = 2
)

type stylishFormatter struct{}

func (f stylishFormatter) Format(root *diff.Node) (string, error) {
	return formatStylish(root.Children, 1), nil
}

func formatStylish(nodes []*diff.Node, depth int) string {
	lines := make([]string, 0, len(nodes)*2)
	lines = append(lines, "{")

	for _, node := range nodes {
		prefix := stylishPrefix(node.Status, depth)

		switch node.Status {
		case diff.StatusRemoved:
			lines = append(lines, formatLine(prefix, "- ", node.Key, formatStylishValue(node.OldValue, depth)))
		case diff.StatusAdded:
			lines = append(lines, formatLine(prefix, "+ ", node.Key, formatStylishValue(node.NewValue, depth)))
		case diff.StatusUpdated:
			removedLine := formatLine(prefix, "- ", node.Key, formatStylishValue(node.OldValue, depth))
			addedLine := formatLine(prefix, "+ ", node.Key, formatStylishValue(node.NewValue, depth))
			lines = append(lines, fmt.Sprintf("%s\n%s", removedLine, addedLine))
		case diff.StatusUnchanged:
			lines = append(lines, formatLine(prefix, "", node.Key, formatStylishValue(node.NewValue, depth)))
		case diff.StatusNested:
			lines = append(lines, formatLine(prefix, "", node.Key, formatStylish(node.Children, depth+1)))
		}
	}

	lines = append(lines, stylishClosingBrace(depth)+"}")

	return strings.Join(lines, "\n")
}

func stylishPrefix(status diff.Status, depth int) string {
	keyColumn := indentStep * depth
	if status == diff.StatusAdded || status == diff.StatusUpdated || status == diff.StatusRemoved {
		return strings.Repeat(" ", keyColumn-markerWidth)
	}

	return strings.Repeat(" ", keyColumn)
}

func stylishClosingBrace(depth int) string {
	return strings.Repeat(" ", indentStep*(depth-1))
}

func formatLine(indent, sign, key, value string) string {
	return fmt.Sprintf("%s%s%s: %s", indent, sign, key, value)
}

func toString(val any) string {
	if val == nil {
		return "null"
	}

	return fmt.Sprintf("%v", val)
}

func formatStylishValue(v any, depth int) string {
	switch value := v.(type) {
	case map[string]any:
		return formatMap(value, depth)
	case []any:
		return formatSlice(value, depth)
	default:
		return toString(v)
	}
}

func formatSlice(items []any, depth int) string {
	childDepth := depth + 1

	lines := make([]string, 0, len(items)+2)
	lines = append(lines, "[")

	itemIndent := strings.Repeat(" ", indentStep*childDepth)
	for _, item := range items {
		lines = append(lines, itemIndent+formatStylishValue(item, childDepth))
	}

	lines = append(lines, stylishClosingBrace(childDepth)+"]")

	return strings.Join(lines, "\n")
}

func formatMap(m map[string]any, depth int) string {
	childDepth := depth + 1

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	lines := make([]string, 0, len(m)+2)
	lines = append(lines, "{")

	keyIndent := strings.Repeat(" ", indentStep*childDepth)
	for _, k := range keys {
		lines = append(lines, formatLine(keyIndent, "", k, formatStylishValue(m[k], childDepth)))
	}

	lines = append(lines, stylishClosingBrace(childDepth)+"}")

	return strings.Join(lines, "\n")
}
