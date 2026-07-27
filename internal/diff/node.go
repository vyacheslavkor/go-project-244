package diff

type Node struct {
	Status   string
	OldValue any
	Value    any
}

func NewNode(status string, oldValue, value any) *Node {
	return &Node{
		Status:   status,
		OldValue: oldValue,
		Value:    value,
	}
}
