package formatters

import (
	"code/internal/diff"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		name       string
		format     string
		beforeFile string
		afterFile  string
		wantFile   string
	}{
		{
			name:       "stylish renders all statuses nested values and scalars",
			format:     "stylish",
			beforeFile: "sample_before.json",
			afterFile:  "sample_after.json",
			wantFile:   "expected_stylish.txt",
		},
		{
			name:       "stylish renders empty tree as braces",
			format:     "stylish",
			beforeFile: "empty.json",
			afterFile:  "empty.json",
			wantFile:   "expected_stylish_empty.txt",
		},
		{
			name:       "plain renders changed properties and skips unchanged",
			format:     "plain",
			beforeFile: "sample_before.json",
			afterFile:  "sample_after.json",
			wantFile:   "expected_plain.txt",
		},
		{
			name:       "plain renders empty tree as empty string",
			format:     "plain",
			beforeFile: "empty.json",
			afterFile:  "empty.json",
			wantFile:   "expected_plain_empty.txt",
		},
		{
			name:       "plain renders unchanged nested tree as empty string",
			format:     "plain",
			beforeFile: "unchanged_nested.json",
			afterFile:  "unchanged_nested.json",
			wantFile:   "expected_plain_empty.txt",
		},
		{
			name:       "json renders root tree with all node statuses",
			format:     "json",
			beforeFile: "sample_before.json",
			afterFile:  "sample_after.json",
			wantFile:   "expected_json.txt",
		},
		{
			name:       "json renders empty tree as root without children",
			format:     "json",
			beforeFile: "empty.json",
			afterFile:  "empty.json",
			wantFile:   "expected_json_empty.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree := buildDiffTree(t, tc.beforeFile, tc.afterFile)

			formatter, err := NewFormatter(tc.format)
			require.NoError(t, err)

			got, err := formatter.Format(tree)
			assert.NoError(t, err)
			assert.Equal(t, readFixture(t, tc.wantFile), got)
		})
	}
}

func buildDiffTree(t *testing.T, beforeFile, afterFile string) *diff.Node {
	t.Helper()

	return diff.Build(loadMapFixture(t, beforeFile), loadMapFixture(t, afterFile))
}

func loadMapFixture(t *testing.T, name string) map[string]any {
	t.Helper()

	var data map[string]any
	require.NoError(t, json.Unmarshal([]byte(readFixture(t, name)), &data))

	return data
}

func readFixture(t *testing.T, name string) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Clean(filepath.Join("testdata", "fixture", name)))
	require.NoError(t, err)

	return strings.TrimSpace(string(content))
}
