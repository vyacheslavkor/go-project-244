package diff

// Status describes how a property changed between two configurations.
type Status string

// Node status values used in a difference tree.
const (
	StatusAdded     Status = "added"
	StatusRemoved   Status = "removed"
	StatusUpdated   Status = "updated"
	StatusNested    Status = "nested"
	StatusUnchanged Status = "unchanged"
	StatusRoot      Status = "root"
)

// Node is a difference tree node. The document root uses StatusRoot
// with an empty Key; property nodes carry Key and optional Children.
type Node struct {
	// Key is the property name; empty on the document root.
	Key string
	// Status is how the property changed.
	Status Status
	// OldValue is the value from the first configuration, if any.
	OldValue any
	// NewValue is the value from the second configuration, if any.
	NewValue any
	// Children holds nested differences for root and nested nodes.
	Children []*Node
}
