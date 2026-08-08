package formatters

import (
	"code/internal/diff"
	"fmt"
	"strings"
)

type plainFormatter struct{}

func (f plainFormatter) Format(root *diff.Node) (string, error) {
	return formatPlain(root.Children, ""), nil
}

func formatPlain(nodes []*diff.Node, path string) string {
	lines := make([]string, 0, len(nodes))

	for _, node := range nodes {
		switch node.Status {
		case diff.StatusAdded:
			lines = append(lines, fmt.Sprintf("Property '%s%s' was added with value: %s", path, node.Key, formatValueToPlain(node.NewValue)))
		case diff.StatusRemoved:
			lines = append(lines, fmt.Sprintf("Property '%s%s' was removed", path, node.Key))
		case diff.StatusUpdated:
			lines = append(lines, fmt.Sprintf("Property '%s%s' was updated. From %s to %s", path, node.Key, formatValueToPlain(node.OldValue), formatValueToPlain(node.NewValue)))
		case diff.StatusUnchanged:
			continue
		case diff.StatusNested:
			nested := formatPlain(node.Children, fmt.Sprintf("%s%s.", path, node.Key))
			if nested != "" {
				lines = append(lines, nested)
			}
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
