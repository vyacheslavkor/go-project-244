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

func NewDiff(parsed1, parsed2 map[string]any) Diff {
	keys := extractKeys(parsed1, parsed2)
	result := Diff{Nodes: map[string]*Node{}}

	for _, k := range keys {
		if _, ok := parsed1[k]; !ok {
			val := parsed2[k]
			result.Nodes[k] = NewNode("added", nil, val)

			continue
		}

		if _, ok := parsed2[k]; !ok {
			val := parsed1[k]
			result.Nodes[k] = NewNode("removed", val, nil)

			continue
		}

		oldVal := parsed1[k]
		val := parsed2[k]

		if reflect.TypeOf(oldVal) == reflect.TypeOf(val) && oldVal == val {
			result.Nodes[k] = NewNode("unchanged", oldVal, val)

			continue
		}

		result.Nodes[k] = NewNode("updated", oldVal, val)
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
