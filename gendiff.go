package code

import (
	"code/internal/diff"
	"code/internal/parsers"
	"fmt"
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

func formatOutput(d diff.Diff) string {
	keys := d.ExtractKeys()
	depth := 1

	lines := make([]string, 0, len(d.Nodes)*2)
	lines = append(lines, "{")

	indentSize := replacerCount*depth + replacerCount*(depth-1)
	bracketIndent := strings.Repeat(replacer, indentSize-replacerCount)

	for _, k := range keys {
		node := d.Nodes[k]
		lines = append(lines, GenerateOutputStr(node, k, 1, replacer, indentSize))
	}

	lines = append(lines, bracketIndent+"}")

	result := strings.Join(lines, "\n")

	return result
}

func GenerateOutputStr(node *diff.Node, key string, depth int, replacer string, indentSize int) string {
	switch node.Status {
	case "removed":
		return fmt.Sprintf("%s- %s: %v", getIndent(node, indentSize, replacerCount), key, node.OldValue)
	case "added":
		return fmt.Sprintf("%s+ %s: %v", getIndent(node, indentSize, replacerCount), key, node.Value)
	case "updated":
		removedLine := fmt.Sprintf("%s- %s: %v", getIndent(node, indentSize, replacerCount), key, node.OldValue)
		addedLine := fmt.Sprintf("%s+ %s: %v", getIndent(node, indentSize, replacerCount), key, node.Value)
		return fmt.Sprintf("%s\n%s", removedLine, addedLine)
	case "unchanged":
		return fmt.Sprintf("%s%s: %v", getIndent(node, indentSize, replacerCount), key, node.Value)
	}

	return ""
}

func getIndent(node *diff.Node, indentSize, replacerCount int) string {
	if node.Status == "added" || node.Status == "updated" || node.Status == "removed" {
		return strings.Repeat(replacer, indentSize)
	}

	return strings.Repeat(replacer, indentSize+replacerCount)
}
