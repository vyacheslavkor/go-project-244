package diff

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	nested1 := map[string]any{
		"keep":   true,
		"change": "old",
		"gone":   float64(1),
	}
	nested2 := map[string]any{
		"keep":   true,
		"change": "new",
		"fresh":  float64(2),
	}

	parsed1 := map[string]any{
		"unchanged_string": "hello",
		"unchanged_bool":   true,
		"unchanged_null":   nil,
		"unchanged_number": float64(42),
		"unchanged_list":   []any{"a", float64(1), false},
		"updated_string":   "before",
		"updated_number":   float64(10),
		"updated_bool":     false,
		"updated_null":     "not-null",
		"updated_list":     []any{"x"},
		"updated_type":     "42",
		"removed":          "bye",
		"nested":           nested1,
		"map_to_scalar": map[string]any{
			"inner": "value",
		},
	}
	parsed2 := map[string]any{
		"unchanged_string": "hello",
		"unchanged_bool":   true,
		"unchanged_null":   nil,
		"unchanged_number": float64(42),
		"unchanged_list":   []any{"a", float64(1), false},
		"updated_string":   "after",
		"updated_number":   float64(20),
		"updated_bool":     true,
		"updated_null":     nil,
		"updated_list":     []any{"x", "y"},
		"updated_type":     float64(42),
		"added":            "hi",
		"nested":           nested2,
		"map_to_scalar":    "plain",
		"scalar_to_map": map[string]any{
			"inner": "value",
		},
	}

	got := Build(parsed1, parsed2)

	expected := &Node{
		Key:      "",
		Status:   StatusRoot,
		OldValue: nil,
		Value:    nil,
		Children: []*Node{
			{Key: "added", Status: StatusAdded, Value: "hi"},
			{Key: "map_to_scalar", Status: StatusUpdated, OldValue: map[string]any{"inner": "value"}, Value: "plain"},
			{Key: "nested", Status: StatusNested, OldValue: nil, Value: nil, Children: []*Node{
				{Key: "change", Status: StatusUpdated, OldValue: "old", Value: "new"},
				{Key: "fresh", Status: StatusAdded, Value: float64(2)},
				{Key: "gone", Status: StatusRemoved, OldValue: float64(1)},
				{Key: "keep", Status: StatusUnchanged, OldValue: true, Value: true},
			}},
			{Key: "removed", Status: StatusRemoved, OldValue: "bye"},
			{Key: "scalar_to_map", Status: StatusAdded, Value: map[string]any{"inner": "value"}},
			{Key: "unchanged_bool", Status: StatusUnchanged, OldValue: true, Value: true},
			{Key: "unchanged_list", Status: StatusUnchanged, OldValue: []any{"a", float64(1), false}, Value: []any{"a", float64(1), false}},
			{Key: "unchanged_null", Status: StatusUnchanged, OldValue: nil, Value: nil},
			{Key: "unchanged_number", Status: StatusUnchanged, OldValue: float64(42), Value: float64(42)},
			{Key: "unchanged_string", Status: StatusUnchanged, OldValue: "hello", Value: "hello"},
			{Key: "updated_bool", Status: StatusUpdated, OldValue: false, Value: true},
			{Key: "updated_list", Status: StatusUpdated, OldValue: []any{"x"}, Value: []any{"x", "y"}},
			{Key: "updated_null", Status: StatusUpdated, OldValue: "not-null", Value: nil},
			{Key: "updated_number", Status: StatusUpdated, OldValue: float64(10), Value: float64(20)},
			{Key: "updated_string", Status: StatusUpdated, OldValue: "before", Value: "after"},
			{Key: "updated_type", Status: StatusUpdated, OldValue: "42", Value: float64(42)},
		},
	}

	require.Equal(t, expected, got)
}
