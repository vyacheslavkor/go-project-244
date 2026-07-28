package formatters

import (
	"code/internal/diff"
	"encoding/json"
)

type JsonFormatter struct{}

type jsonNode struct {
	Key      string     `json:"key"`
	Status   string     `json:"status"`
	OldValue *any       `json:"old_value,omitempty"`
	Value    *any       `json:"value,omitempty"`
	Children []jsonNode `json:"children,omitempty"`
}

func (f *JsonFormatter) Format(d *diff.Diff) (string, error) {
	nodes := processDiffJson(d)
	rootNode := jsonNode{
		Key:    "",
		Status: "root",
		Children: nodes,
	}

	data, err := json.Marshal(rootNode)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func processDiffJson(d *diff.Diff) []jsonNode {
	keys := d.ExtractKeys()

	nodes := make([]jsonNode, 0, len(keys))

	for _, k := range keys {
		node := d.Nodes[k]
		n := jsonNode{
			Key:    k,
			Status: node.Status,
		}

		switch node.Status {
		case "added":
			n.Value = &node.Value
		case "removed":
			n.OldValue = &node.OldValue
		case "updated":
			n.Value = &node.Value
			n.OldValue = &node.OldValue
		case "unchanged":
			n.Value = &node.Value
		case "nested":
			n.Children = processDiffJson(node.Children)
		}

		nodes = append(nodes, n)
	}

	return nodes
}
