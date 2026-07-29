package formatters

import (
	"code/internal/diff"
	"encoding/json"
)

// jsonStatusRoot is a wire-only status for the JSON output root object.
// It is not a [diff.Status]: the internal difference tree has no root node,
// while the JSON format wraps children in a synthetic envelope.
const jsonStatusRoot = "root"

type jsonFormatter struct{}

// JSON output is an object with fields:
//
//	key       string           // property name; empty on the synthetic root
//	status    string           // "root" | added | removed | updated | nested | unchanged
//	old_value any, optional    // present for removed/updated
//	value     any, optional    // present for added/updated/unchanged
//	children  []object, optional // present for root/nested
//
// Status "root" exists only in this JSON envelope. Tree node statuses
// ([diff.Status]) never include "root"; the overlap of the other names is
// intentional but the two layers are separate contracts.
type jsonNode struct {
	Key      string     `json:"key"`
	Status   string     `json:"status"`
	OldValue *any       `json:"old_value,omitempty"`
	Value    *any       `json:"value,omitempty"`
	Children []jsonNode `json:"children,omitempty"`
}

func (f jsonFormatter) Format(t *diff.Tree) (string, error) {
	nodes := buildJSONNodes(t)
	rootNode := jsonNode{
		Key:      "",
		Status:   jsonStatusRoot,
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

		nodes = append(nodes, n)
	}

	return nodes
}
