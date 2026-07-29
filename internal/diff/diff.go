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
			result.Nodes[k] = &Node{Status: StatusAdded, Value: val2}

			continue
		}

		if !has2 {
			result.Nodes[k] = &Node{Status: StatusRemoved, OldValue: val1}

			continue
		}

		map1, isMap1 := val1.(map[string]any)
		map2, isMap2 := val2.(map[string]any)

		if isMap1 && isMap2 {
			result.Nodes[k] = &Node{
				Status:   StatusNested,
				OldValue: val1,
				Value:    val2,
				Children: NewTree(map1, map2),
			}

			continue
		}

		if reflect.DeepEqual(val1, val2) {
			result.Nodes[k] = &Node{Status: StatusUnchanged, OldValue: val1, Value: val2}

			continue
		}

		result.Nodes[k] = &Node{Status: StatusUpdated, OldValue: val1, Value: val2}
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
