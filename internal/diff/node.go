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
)

// Node is a single property difference in a [Tree].
type Node struct {
	// Status is how the property changed.
	Status Status
	// OldValue is the value from the first configuration, if any.
	OldValue any
	// Value is the value from the second configuration, if any.
	Value any
	// Children holds nested differences when Status is StatusNested.
	Children *Tree
}
