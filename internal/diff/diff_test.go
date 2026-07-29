package diff

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTree(t *testing.T) {
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

	got := NewTree(parsed1, parsed2)

	expected := &Tree{
		Nodes: map[string]*Node{
			"added": {
				Status: StatusAdded,
				Value:  "hi",
			},
			"map_to_scalar": {
				Status:   StatusUpdated,
				OldValue: map[string]any{"inner": "value"},
				Value:    "plain",
			},
			"nested": {
				Status:   StatusNested,
				OldValue: nested1,
				Value:    nested2,
				Children: &Tree{
					Nodes: map[string]*Node{
						"change": {
							Status:   StatusUpdated,
							OldValue: "old",
							Value:    "new",
						},
						"fresh": {
							Status: StatusAdded,
							Value:  float64(2),
						},
						"gone": {
							Status:   StatusRemoved,
							OldValue: float64(1),
						},
						"keep": {
							Status:   StatusUnchanged,
							OldValue: true,
							Value:    true,
						},
					},
				},
			},
			"removed": {
				Status:   StatusRemoved,
				OldValue: "bye",
			},
			"scalar_to_map": {
				Status: StatusAdded,
				Value: map[string]any{
					"inner": "value",
				},
			},
			"unchanged_bool": {
				Status:   StatusUnchanged,
				OldValue: true,
				Value:    true,
			},
			"unchanged_list": {
				Status:   StatusUnchanged,
				OldValue: []any{"a", float64(1), false},
				Value:    []any{"a", float64(1), false},
			},
			"unchanged_null": {
				Status:   StatusUnchanged,
				OldValue: nil,
				Value:    nil,
			},
			"unchanged_number": {
				Status:   StatusUnchanged,
				OldValue: float64(42),
				Value:    float64(42),
			},
			"unchanged_string": {
				Status:   StatusUnchanged,
				OldValue: "hello",
				Value:    "hello",
			},
			"updated_bool": {
				Status:   StatusUpdated,
				OldValue: false,
				Value:    true,
			},
			"updated_list": {
				Status:   StatusUpdated,
				OldValue: []any{"x"},
				Value:    []any{"x", "y"},
			},
			"updated_null": {
				Status:   StatusUpdated,
				OldValue: "not-null",
				Value:    nil,
			},
			"updated_number": {
				Status:   StatusUpdated,
				OldValue: float64(10),
				Value:    float64(20),
			},
			"updated_string": {
				Status:   StatusUpdated,
				OldValue: "before",
				Value:    "after",
			},
			"updated_type": {
				Status:   StatusUpdated,
				OldValue: "42",
				Value:    float64(42),
			},
		},
	}

	require.Equal(t, expected, got)
}
