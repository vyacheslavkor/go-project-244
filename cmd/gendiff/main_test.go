package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppExecution(t *testing.T) {
	compiledPath := buildBinary(t)
	ctx := context.Background()
	help := readHelpFixture(t)

	t.Run("prints help and exits zero", func(t *testing.T) {
		stdout, stderr, err := runBinary(t, ctx, compiledPath, "-h")
		require.NoError(t, err)
		assert.Equal(t, help, strings.TrimSpace(stdout))
		assert.Empty(t, stderr)
	})

	t.Run("prints help and usage reason on missing args", func(t *testing.T) {
		stdout, stderr, err := runBinary(t, ctx, compiledPath)
		require.Error(t, err)
		assert.Equal(t, 1, exitCode(t, err))
		assert.Equal(t, help, strings.TrimSpace(stdout))
		assert.Equal(t, "incorrect usage: expected 2 arguments, got 0\n", stderr)
	})

	t.Run("prints only reason for missing files", func(t *testing.T) {
		stdout, stderr, err := runBinary(t, ctx, compiledPath, "missing1.json", "missing2.json")
		require.Error(t, err)
		assert.Equal(t, 1, exitCode(t, err))
		assert.Empty(t, stdout)
		assert.Equal(
			t,
			"failed to read file: 'missing1.json': stat missing1.json: no such file or directory\n",
			stderr,
		)
	})
}

func buildBinary(t *testing.T) string {
	t.Helper()

	compiledPath := filepath.Join(t.TempDir(), "gendiff")
	ctx := context.Background()

	//nolint:gosec // G204: paths are controlled by the test, no risk of command injection
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", compiledPath, ".")
	require.NoError(t, buildCmd.Run(), "failed to compile binary")

	return compiledPath
}

func runBinary(t *testing.T, ctx context.Context, compiledPath string, args ...string) (stdout, stderr string, err error) {
	t.Helper()

	//nolint:gosec // G204: paths are controlled by the test, no risk of command injection
	cmd := exec.CommandContext(ctx, compiledPath, args...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	runErr := cmd.Run()

	return outBuf.String(), errBuf.String(), runErr
}

func exitCode(t *testing.T, err error) int {
	t.Helper()

	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr)

	return exitErr.ExitCode()
}

func readHelpFixture(t *testing.T) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Clean(filepath.Join(
		"..", "..", "internal", "cli", "testdata", "fixture", "help.txt",
	)))
	require.NoError(t, err)

	return strings.TrimSpace(string(content))
}
