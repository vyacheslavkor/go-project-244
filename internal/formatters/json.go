package formatters

import (
	"code/internal/diff"
	"encoding/json"
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
		return "", err
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
			n.Value = &node.Value
		case diff.StatusRemoved:
			n.OldValue = &node.OldValue
		case diff.StatusUpdated:
			n.Value = &node.Value
			n.OldValue = &node.OldValue
		case diff.StatusUnchanged:
			n.Value = &node.Value
		case diff.StatusNested:
			n.Children = buildJSONNodes(node.Children)
		}

		result = append(result, n)
	}

	return result
}
