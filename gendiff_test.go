package code

import (
	"code/internal/formatters"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenDiff(t *testing.T) {
	dir := t.TempDir()
	file1 := writeTempFile(t, dir, "a.json", `{"k":"v1"}`)
	file2 := writeTempFile(t, dir, "b.json", `{"k":"v2"}`)

	t.Run("returns stylish difference for matching json files", func(t *testing.T) {
		got, err := GenDiff(file1, file2, "stylish")
		assert.NoError(t, err)
		assert.Equal(t, "{\n  - k: v1\n  + k: v2\n}", got)
	})

	t.Run("returns plain difference when plain format is requested", func(t *testing.T) {
		got, err := GenDiff(file1, file2, "plain")
		assert.NoError(t, err)
		assert.Equal(t, "Property 'k' was updated. From 'v1' to 'v2'", got)
	})

	t.Run("rejects unknown output format", func(t *testing.T) {
		got, err := GenDiff(file1, file2, "xml")
		assert.ErrorIs(t, err, formatters.ErrInvalidFormat)
		assert.Empty(t, got)
	})

	t.Run("rejects missing input file", func(t *testing.T) {
		got, err := GenDiff(file1, filepath.Join(dir, "missing.json"), "stylish")
		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.Empty(t, got)
	})
}

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	return path
}
