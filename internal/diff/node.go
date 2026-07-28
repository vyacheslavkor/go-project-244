package diff

type Node struct {
	Status   string
	OldValue any
	Value    any
	Children *Diff
}

func NewNode(status string, oldValue, value any, children *Diff) *Node {
	return &Node{
		Status:   status,
		OldValue: oldValue,
		Value:    value,
		Children: children,
	}
}
