package diff

import (
	"reflect"
	"slices"
)

// Build builds a difference tree by comparing before and after.
// The returned node is the document root (StatusRoot) whose Children
// hold the property differences.
func Build(before, after map[string]any) *Node {
	return &Node{Status: StatusRoot, Children: buildChildren(before, after)}
}

func buildChildren(before, after map[string]any) []*Node {
	keys := sortedUnionKeys(before, after)
	children := make([]*Node, 0, len(keys))

	for _, k := range keys {
		beforeValue, existsBefore := before[k]
		afterValue, existsAfter := after[k]

		if !existsBefore {
			children = append(children, &Node{Key: k, Status: StatusAdded, Value: afterValue})

			continue
		}

		if !existsAfter {
			children = append(children, &Node{Key: k, Status: StatusRemoved, OldValue: beforeValue})

			continue
		}

		beforeMap, beforeIsMap := beforeValue.(map[string]any)
		afterMap, afterIsMap := afterValue.(map[string]any)

		if beforeIsMap && afterIsMap {
			children = append(children, &Node{
				Key:      k,
				Status:   StatusNested,
				Children: buildChildren(beforeMap, afterMap),
			})

			continue
		}

		if reflect.DeepEqual(beforeValue, afterValue) {
			children = append(children, &Node{Key: k, Status: StatusUnchanged, OldValue: beforeValue, Value: afterValue})

			continue
		}

		children = append(children, &Node{Key: k, Status: StatusUpdated, OldValue: beforeValue, Value: afterValue})
	}

	return children
}

func sortedUnionKeys(a, b map[string]any) []string {
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
