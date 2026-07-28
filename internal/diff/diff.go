package diff

import (
	"reflect"
	"slices"
)

type Diff struct {
	Nodes map[string]*Node
}

func (d *Diff) ExtractKeys() []string {
	result := make([]string, 0, len(d.Nodes))

	for k := range d.Nodes {
		result = append(result, k)
	}

	slices.Sort(result)

	return result
}

func NewDiff(parsed1, parsed2 map[string]any) *Diff {
	keys := extractKeys(parsed1, parsed2)
	result := &Diff{Nodes: map[string]*Node{}}

	for _, k := range keys {
		val1, has1 := parsed1[k]
		val2, has2 := parsed2[k]

		if !has1 {
			result.Nodes[k] = NewNode("added", nil, val2, nil)

			continue
		}

		if !has2 {
			result.Nodes[k] = NewNode("removed", val1, nil, nil)

			continue
		}

		map1, isMap1 := val1.(map[string]any)
		map2, isMap2 := val2.(map[string]any)

		if isMap1 && isMap2 {
			diff := NewDiff(map1, map2)
			result.Nodes[k] = NewNode("nested", val1, val2, diff)

			continue
		}

		if reflect.DeepEqual(val1, val2) {
			result.Nodes[k] = NewNode("unchanged", val1, val2, nil)

			continue
		}

		result.Nodes[k] = NewNode("updated", val1, val2, nil)
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
