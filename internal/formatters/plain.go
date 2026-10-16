package formatters

import (
	"code/internal/diff"
	"fmt"
	"strconv"
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
		default:
			panic(fmt.Sprintf("plain: unknown status %q", node.Status))
		}
	}

	return strings.Join(lines, "\n")
}

func formatValueToPlain(val any) string {
	switch v := val.(type) {
	case nil:
		return "null"
	case string:
		return fmt.Sprintf("'%s'", v)
	case bool:
		return strconv.FormatBool(v)
	case float64:
		return formatNumber(v)
	case map[string]any, []any:
		return "[complex value]"
	default:
		return fmt.Sprintf("%v", v)
	}
}
