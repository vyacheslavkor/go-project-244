package formatters

import (
	"code/internal/diff"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	sampleBefore = map[string]any{
		"added_later": "x",
		"removed":     "bye",
		"updated":     "before",
		"unchanged":   true,
		"null_val":    "not-null",
		"number":      float64(10),
		"list":        []any{"a"},
		"complex":     map[string]any{"inner": "v"},
		"nested": map[string]any{
			"keep":   true,
			"change": "old",
			"gone":   float64(1),
		},
	}
	sampleAfter = map[string]any{
		"added":     "hi",
		"updated":   "after",
		"unchanged": true,
		"null_val":  nil,
		"number":    float64(20),
		"list":      []any{"a", "b"},
		"complex":   "plain",
		"nested": map[string]any{
			"keep":   true,
			"change": "new",
			"fresh":  float64(2),
		},
		"bool_add": false,
	}
	sampleDiffTree          = diff.NewTree(sampleBefore, sampleAfter)
	emptyDiffTree           = &diff.Tree{Nodes: map[string]*diff.Node{}}
	unchangedNestedDiffTree = diff.NewTree(
		map[string]any{"nested": map[string]any{"key": "value"}},
		map[string]any{"nested": map[string]any{"key": "value"}},
	)
)

func TestNewFormatter(t *testing.T) {
	testCases := []struct {
		name    string
		format  string
		wantErr error
	}{
		{name: "creates stylish formatter", format: "stylish"},
		{name: "creates plain formatter", format: "plain"},
		{name: "creates json formatter", format: "json"},
		{name: "rejects unknown format", format: "xml", wantErr: ErrInvalidFormat},
		{name: "treats empty format as stylish", format: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewFormatter(tc.format)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}

func TestFormatters(t *testing.T) {
	testCases := []struct {
		name     string
		format   string
		tree     *diff.Tree
		wantFile string
	}{
		{
			name:     "stylish renders all statuses nested values and scalars",
			format:   "stylish",
			tree:     sampleDiffTree,
			wantFile: "expected_stylish.txt",
		},
		{
			name:     "stylish renders empty tree as braces",
			format:   "stylish",
			tree:     emptyDiffTree,
			wantFile: "expected_stylish_empty.txt",
		},
		{
			name:     "plain renders changed properties and skips unchanged",
			format:   "plain",
			tree:     sampleDiffTree,
			wantFile: "expected_plain.txt",
		},
		{
			name:     "plain renders empty tree as empty string",
			format:   "plain",
			tree:     emptyDiffTree,
			wantFile: "expected_plain_empty.txt",
		},
		{
			name:     "plain renders unchanged nested tree as empty string",
			format:   "plain",
			tree:     unchangedNestedDiffTree,
			wantFile: "expected_plain_empty.txt",
		},
		{
			name:     "json renders root tree with all node statuses",
			format:   "json",
			tree:     sampleDiffTree,
			wantFile: "expected_json.txt",
		},
		{
			name:     "json renders empty tree as root without children",
			format:   "json",
			tree:     emptyDiffTree,
			wantFile: "expected_json_empty.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formatter, err := NewFormatter(tc.format)
			require.NoError(t, err)

			got, err := formatter.Format(tc.tree)
			assert.NoError(t, err)
			assert.Equal(t, readFixture(t, tc.wantFile), got)
		})
	}
}

func readFixture(t *testing.T, name string) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Clean(filepath.Join("testdata", "fixture", name)))
	require.NoError(t, err)

	return strings.TrimSpace(string(content))
}
