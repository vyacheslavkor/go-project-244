package formatters

import (
	"code/internal/diff"
	"fmt"
	"slices"
	"strings"
)

const (
	replacer      = " "
	replacerCount = 2
)

type stylishFormatter struct{}

func (f stylishFormatter) Format(t *diff.Tree) (string, error) {
	return formatStylish(t, 1), nil
}

func formatStylish(t *diff.Tree, depth int) string {
	indentSize := replacerCount*depth + replacerCount*(depth-1)
	bracketIndent := strings.Repeat(replacer, indentSize-replacerCount)

	keys := t.ExtractKeys()

	lines := make([]string, 0, len(t.Nodes)*2)
	lines = append(lines, "{")

	for _, k := range keys {
		node := t.Nodes[k]
		indent := getIndent(node, indentSize, replacerCount)

		switch node.Status {
		case diff.StatusRemoved:
			lines = append(lines, formatLine(indent, "- ", k, formatStylishValue(node.OldValue, depth)))
		case diff.StatusAdded:
			lines = append(lines, formatLine(indent, "+ ", k, formatStylishValue(node.Value, depth)))
		case diff.StatusUpdated:
			removedLine := formatLine(indent, "- ", k, formatStylishValue(node.OldValue, depth))
			addedLine := formatLine(indent, "+ ", k, formatStylishValue(node.Value, depth))
			lines = append(lines, fmt.Sprintf("%s\n%s", removedLine, addedLine))
		case diff.StatusUnchanged:
			lines = append(lines, formatLine(indent, "", k, formatStylishValue(node.Value, depth)))
		case diff.StatusNested:
			lines = append(lines, formatLine(indent, "", k, formatStylish(node.Children, depth+1)))
		}
	}

	lines = append(lines, bracketIndent+"}")

	return strings.Join(lines, "\n")
}

func getIndent(node *diff.Node, indentSize, replacerCount int) string {
	if node.Status == diff.StatusAdded || node.Status == diff.StatusUpdated || node.Status == diff.StatusRemoved {
		return strings.Repeat(replacer, indentSize)
	}

	return strings.Repeat(replacer, indentSize+replacerCount)
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
	if m, isMap := v.(map[string]any); isMap {
		return formatMap(m, depth)
	}

	return toString(v)
}

func formatMap(m map[string]any, depth int) string {
	childDepth := depth + 1
	childIndentSize := replacerCount*childDepth + replacerCount*(childDepth-1)
	keyIndent := strings.Repeat(replacer, childIndentSize+replacerCount)
	bracketIndent := strings.Repeat(replacer, childIndentSize-replacerCount)

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	lines := make([]string, 0, len(m)+2)
	lines = append(lines, "{")

	for _, k := range keys {
		lines = append(lines, formatLine(keyIndent, "", k, formatStylishValue(m[k], childDepth)))
	}

	lines = append(lines, bracketIndent+"}")

	return strings.Join(lines, "\n")
}
