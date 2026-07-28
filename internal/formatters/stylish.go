package formatters

import (
	"code/internal/diff"
	"fmt"
	"sort"
	"strings"
)

const (
	replacer      = " "
	replacerCount = 2
)

type stylishFormatter struct{}

func (f stylishFormatter) Format(d *diff.Diff) (string, error) {
	return processDiffStylish(d, 1), nil
}

func processDiffStylish(d *diff.Diff, depth int) string {
	indentSize := replacerCount*depth + replacerCount*(depth-1)
	bracketIndent := strings.Repeat(replacer, indentSize-replacerCount)

	keys := d.ExtractKeys()

	lines := make([]string, 0, len(d.Nodes)*2)
	lines = append(lines, "{")

	for _, k := range keys {
		node := d.Nodes[k]
		indent := getIndent(node, indentSize, replacerCount)

		switch node.Status {
		case "removed":
			lines = append(lines, formatLine(indent, "- ", k, formatStylishValue(node.OldValue, depth)))
		case "added":
			lines = append(lines, formatLine(indent, "+ ", k, formatStylishValue(node.Value, depth)))
		case "updated":
			removedLine := formatLine(indent, "- ", k, formatStylishValue(node.OldValue, depth))
			addedLine := formatLine(indent, "+ ", k, formatStylishValue(node.Value, depth))
			lines = append(lines, fmt.Sprintf("%s\n%s", removedLine, addedLine))
		case "unchanged":
			lines = append(lines, formatLine(indent, "", k, formatStylishValue(node.Value, depth)))
		case "nested":
			lines = append(lines, formatLine(indent, "", k, processDiffStylish(node.Children, depth+1)))
		}
	}

	lines = append(lines, bracketIndent+"}")

	return strings.Join(lines, "\n")
}

func getIndent(node *diff.Node, indentSize, replacerCount int) string {
	if node.Status == "added" || node.Status == "updated" || node.Status == "removed" {
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
	sort.Strings(keys)

	lines := make([]string, 0, len(m)+2)
	lines = append(lines, "{")

	for _, k := range keys {
		lines = append(lines, formatLine(keyIndent, "", k, formatStylishValue(m[k], childDepth)))
	}

	lines = append(lines, bracketIndent+"}")

	return strings.Join(lines, "\n")
}
