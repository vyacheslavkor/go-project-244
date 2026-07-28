package code

import (
	"code/internal/diff"
	"code/internal/parsers"
	"fmt"
	"sort"
	"strings"
)

func GenDiff(filepath1, filepath2, format string) (string, error) {
	parsed1, parsed2, err := parsers.ParseFiles(filepath1, filepath2)
	if err != nil {
		return "", err
	}

	d := diff.NewDiff(parsed1, parsed2)
	formatted := formatOutput(d)

	return formatted, nil
}

const (
	replacer      = " "
	replacerCount = 2
)

func formatOutput(d *diff.Diff) string {
	return stylish(d)
}

func getIndent(node *diff.Node, indentSize, replacerCount int) string {
	if node.Status == "added" || node.Status == "updated" || node.Status == "removed" {
		return strings.Repeat(replacer, indentSize)
	}

	return strings.Repeat(replacer, indentSize+replacerCount)
}

func stylish(d *diff.Diff) string {
	return iter(d, 1)
}

func iter(d *diff.Diff, depth int) string {
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
			lines = append(lines, formatLine(indent, "- ", k, formatValue(node.OldValue, depth)))
		case "added":
			lines = append(lines, formatLine(indent, "+ ", k, formatValue(node.Value, depth)))
		case "updated":
			removedLine := formatLine(indent, "- ", k, formatValue(node.OldValue, depth))
			addedLine := formatLine(indent, "+ ", k, formatValue(node.Value, depth))
			lines = append(lines, fmt.Sprintf("%s\n%s", removedLine, addedLine))
		case "unchanged":
			lines = append(lines, formatLine(indent, "", k, formatValue(node.Value, depth)))
		case "nested":
			lines = append(lines, formatLine(indent, "", k, iter(node.Children, depth+1)))
		}
	}

	lines = append(lines, bracketIndent+"}")

	return strings.Join(lines, "\n")
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

func formatValue(v any, depth int) string {
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
		lines = append(lines, formatLine(keyIndent, "", k, formatValue(m[k], childDepth)))
	}

	lines = append(lines, bracketIndent+"}")

	return strings.Join(lines, "\n")
}
