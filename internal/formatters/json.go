package formatters

import (
	"code/internal/diff"
	"encoding/json"
	"fmt"
)

type jsonFormatter struct{}

// JSON output is a compact single-line object with fields:
//
//	key       string           // property name; empty on the root
//	status    string           // root | added | removed | updated | nested | unchanged
//	old_value any, optional    // present for removed/updated
//	value     any, optional    // present for added/updated/unchanged
//	children  []object, optional // present for root/nested; omitted when empty
type jsonNode struct {
	Key      string     `json:"key"`
	Status   string     `json:"status"`
	OldValue *any       `json:"old_value,omitempty"`
	Value    *any       `json:"value,omitempty"`
	Children []jsonNode `json:"children,omitempty"`
}

func (f jsonFormatter) Format(root *diff.Node) (string, error) {
	rootNode := jsonNode{
		Key:      root.Key,
		Status:   string(root.Status),
		Children: buildJSONNodes(root.Children),
	}

	data, err := json.Marshal(rootNode)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFormatOutput, err)
	}

	return string(data), nil
}

func buildJSONNodes(nodes []*diff.Node) []jsonNode {
	result := make([]jsonNode, 0, len(nodes))

	for _, node := range nodes {
		n := jsonNode{
			Key:    node.Key,
			Status: string(node.Status),
		}

		switch node.Status {
		case diff.StatusAdded:
			formattedValue := formatValueToJSON(node.NewValue)
			n.Value = &formattedValue
		case diff.StatusRemoved:
			formattedOldValue := formatValueToJSON(node.OldValue)
			n.OldValue = &formattedOldValue
		case diff.StatusUpdated:
			formattedValue := formatValueToJSON(node.NewValue)
			n.Value = &formattedValue

			formattedOldValue := formatValueToJSON(node.OldValue)
			n.OldValue = &formattedOldValue
		case diff.StatusUnchanged:
			formattedValue := formatValueToJSON(node.NewValue)
			n.Value = &formattedValue
		case diff.StatusNested:
			n.Children = buildJSONNodes(node.Children)
		}

		result = append(result, n)
	}

	return result
}

func formatValueToJSON(val any) any {
	switch v := val.(type) {
	case float64:
		return json.Number(formatNumber(v))
	case map[string]any:
		return formatJSONMap(v)
	case []any:
		return formatJSONSlice(v)
	default:
		return val
	}
}

func formatJSONSlice(items []any) []any {
	result := make([]any, 0, len(items))
	for _, item := range items {
		result = append(result, formatValueToJSON(item))
	}

	return result
}

func formatJSONMap(items map[string]any) map[string]any {
	result := make(map[string]any, len(items))
	for key, value := range items {
		result[key] = formatValueToJSON(value)
	}

	return result
}
