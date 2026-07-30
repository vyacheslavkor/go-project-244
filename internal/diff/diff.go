package diff

import (
	"reflect"
	"slices"
)

// Tree is a sorted map of difference nodes keyed by property name.
type Tree struct {
	// Nodes maps property names to their difference nodes.
	Nodes map[string]*Node
}

// Keys returns the node keys in lexicographic order.
func (t *Tree) Keys() []string {
	result := make([]string, 0, len(t.Nodes))

	for k := range t.Nodes {
		result = append(result, k)
	}

	slices.Sort(result)

	return result
}

// NewTree builds a difference tree by comparing before and after.
func NewTree(before, after map[string]any) *Tree {
	keys := extractKeys(before, after)
	result := &Tree{Nodes: map[string]*Node{}}

	for _, k := range keys {
		beforeValue, existsBefore := before[k]
		afterValue, existsAfter := after[k]

		if !existsBefore {
			result.Nodes[k] = &Node{Status: StatusAdded, Value: afterValue}

			continue
		}

		if !existsAfter {
			result.Nodes[k] = &Node{Status: StatusRemoved, OldValue: beforeValue}

			continue
		}

		beforeMap, beforeIsMap := beforeValue.(map[string]any)
		afterMap, afterIsMap := afterValue.(map[string]any)

		if beforeIsMap && afterIsMap {
			result.Nodes[k] = &Node{
				Status:   StatusNested,
				OldValue: beforeValue,
				Value:    afterValue,
				Children: NewTree(beforeMap, afterMap),
			}

			continue
		}

		if reflect.DeepEqual(beforeValue, afterValue) {
			result.Nodes[k] = &Node{Status: StatusUnchanged, OldValue: beforeValue, Value: afterValue}

			continue
		}

		result.Nodes[k] = &Node{Status: StatusUpdated, OldValue: beforeValue, Value: afterValue}
	}

	return result
}

func extractKeys(a, b map[string]any) []string {
	seen := make(map[string]struct{}, len(a)+len(b))

	for k := range a {
		seen[k] = struct{}{}
	}

	for k := range b {
		seen[k] = struct{}{}
	}

	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	return keys
}
