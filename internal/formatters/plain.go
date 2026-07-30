package formatters

import (
	"code/internal/diff"
	"fmt"
	"strings"
)

type plainFormatter struct{}

func (f plainFormatter) Format(t *diff.Tree) (string, error) {
	return formatPlain(t, ""), nil
}

func formatPlain(t *diff.Tree, path string) string {
	keys := t.Keys()

	lines := make([]string, 0, len(keys))

	for _, k := range keys {
		node := t.Nodes[k]

		switch node.Status {
		case diff.StatusAdded:
			lines = append(lines, fmt.Sprintf("Property '%s%s' was added with value: %s", path, k, formatValueToPlain(node.Value)))
		case diff.StatusRemoved:
			lines = append(lines, fmt.Sprintf("Property '%s%s' was removed", path, k))
		case diff.StatusUpdated:
			lines = append(lines, fmt.Sprintf("Property '%s%s' was updated. From %s to %s", path, k, formatValueToPlain(node.OldValue), formatValueToPlain(node.Value)))
		case diff.StatusUnchanged:
			continue
		case diff.StatusNested:
			lines = append(lines, formatPlain(node.Children, fmt.Sprintf("%s%s.", path, k)))
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
