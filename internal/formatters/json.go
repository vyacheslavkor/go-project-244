package formatters

import (
	"code/internal/diff"
	"encoding/json"
)

type jsonFormatter struct{}

type jsonNode struct {
	Key      string      `json:"key"`
	Status   diff.Status `json:"status"`
	OldValue *any        `json:"old_value,omitempty"`
	Value    *any        `json:"value,omitempty"`
	Children []jsonNode  `json:"children,omitempty"`
}

func (f jsonFormatter) Format(t *diff.Tree) (string, error) {
	nodes := buildJSONNodes(t)
	rootNode := jsonNode{
		Key:      "",
		Status:   "root",
		Children: nodes,
	}

	data, err := json.Marshal(rootNode)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func buildJSONNodes(t *diff.Tree) []jsonNode {
	keys := t.ExtractKeys()

	nodes := make([]jsonNode, 0, len(keys))

	for _, k := range keys {
		node := t.Nodes[k]
		n := jsonNode{
			Key:    k,
			Status: node.Status,
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

		nodes = append(nodes, n)
	}

	return nodes
}
