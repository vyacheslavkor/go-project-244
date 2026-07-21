package main

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppExecution(t *testing.T) {
	tempDir := t.TempDir()
	compiledPath := filepath.Join(tempDir, "gendiff")

	ctx := context.Background()

	//nolint:gosec // G204: paths are controlled by the test, no risk of command injection
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", compiledPath, ".")

	err := buildCmd.Run()
	require.NoError(t, err, "failed to compile binary")

	//nolint:gosec // G204: paths are controlled by the test, no risk of command injection
	runCmd := exec.CommandContext(ctx, compiledPath, "-h")

	var stdout, stderr bytes.Buffer
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr

	runErr := runCmd.Run()
	require.NoError(t, runErr, "expected command to succeed with exit code 0")
	assert.NotEmpty(t, stdout.String())
	assert.Empty(t, stderr.String())
}
