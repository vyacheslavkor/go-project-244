package formatters

import (
	"code/internal/diff"
	"fmt"
	"strconv"
	"strings"
)

const (
	plainMsgAdded   = "Property '%s%s' was added with value: %s"
	plainMsgRemoved = "Property '%s%s' was removed"
	plainMsgUpdated = "Property '%s%s' was updated. From %s to %s"
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
			formattedValue := formatValueToPlain(node.NewValue)
			lines = append(lines, fmt.Sprintf(plainMsgAdded, path, node.Key, formattedValue))
		case diff.StatusRemoved:
			lines = append(lines, fmt.Sprintf(plainMsgRemoved, path, node.Key))
		case diff.StatusUpdated:
			formattedOldValue := formatValueToPlain(node.OldValue)
			formattedNewValue := formatValueToPlain(node.NewValue)
			lines = append(
				lines,
				fmt.Sprintf(plainMsgUpdated, path, node.Key, formattedOldValue, formattedNewValue),
			)
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
