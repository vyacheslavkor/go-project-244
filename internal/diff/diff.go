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

// ExtractKeys returns the node keys in lexicographic order.
func (t *Tree) ExtractKeys() []string {
	result := make([]string, 0, len(t.Nodes))

	for k := range t.Nodes {
		result = append(result, k)
	}

	slices.Sort(result)

	return result
}

// NewTree builds a difference tree by comparing parsed1 and parsed2.
func NewTree(parsed1, parsed2 map[string]any) *Tree {
	keys := extractKeys(parsed1, parsed2)
	result := &Tree{Nodes: map[string]*Node{}}

	for _, k := range keys {
		val1, has1 := parsed1[k]
		val2, has2 := parsed2[k]

		if !has1 {
			result.Nodes[k] = NewNode(StatusAdded, nil, val2, nil)

			continue
		}

		if !has2 {
			result.Nodes[k] = NewNode(StatusRemoved, val1, nil, nil)

			continue
		}

		map1, isMap1 := val1.(map[string]any)
		map2, isMap2 := val2.(map[string]any)

		if isMap1 && isMap2 {
			diff := NewTree(map1, map2)
			result.Nodes[k] = NewNode(StatusNested, val1, val2, diff)

			continue
		}

		if reflect.DeepEqual(val1, val2) {
			result.Nodes[k] = NewNode(StatusUnchanged, val1, val2, nil)

			continue
		}

		result.Nodes[k] = NewNode(StatusUpdated, val1, val2, nil)
	}

	return result
}

func extractKeys(a, b map[string]any) []string {
	maxLen := len(a) + len(b)
	uniqueMap := make(map[string]string, maxLen)

	for k := range a {
		_, exists := uniqueMap[k]
		if !exists {
			uniqueMap[k] = k
		}
	}

	for k := range b {
		_, exists := uniqueMap[k]
		if !exists {
			uniqueMap[k] = k
		}
	}

	keys := make([]string, 0, len(uniqueMap))
	for k := range uniqueMap {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	return keys
}
