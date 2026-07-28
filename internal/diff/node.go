package diff

type Status string

const (
	StatusAdded     Status = "added"
	StatusRemoved   Status = "removed"
	StatusUpdated   Status = "updated"
	StatusNested    Status = "nested"
	StatusUnchanged Status = "unchanged"
)

type Node struct {
	Status   Status
	OldValue any
	Value    any
	Children *Tree
}

func NewNode(status Status, oldValue, value any, children *Tree) *Node {
	return &Node{
		Status:   status,
		OldValue: oldValue,
		Value:    value,
		Children: children,
	}
}
