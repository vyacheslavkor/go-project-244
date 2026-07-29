package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	urfave "github.com/urfave/cli/v3"
)

func TestNewCommand(t *testing.T) {
	fixtureDir := getFixturePath()
	json1 := filepath.Join(fixtureDir, "file1.json")
	json2 := filepath.Join(fixtureDir, "file2.json")
	yml1 := filepath.Join(fixtureDir, "file1.yml")
	yml2 := filepath.Join(fixtureDir, "file2.yml")

	expectedStylish := readFixture(t, "expected_stylish.txt")
	expectedPlain := readFixture(t, "expected_plain.txt")
	expectedJSON := readFixture(t, "expected_json.txt")
	help := readFixture(t, "help.txt")

	successCases := []struct {
		name       string
		args       []string
		wantStdout string
	}{
		{
			name:       "prints help when help flag is passed",
			args:       []string{"--help"},
			wantStdout: help,
		},
		{
			name:       "compares json files with default stylish format",
			args:       []string{json1, json2},
			wantStdout: expectedStylish,
		},
		{
			name:       "compares yml files with default stylish format",
			args:       []string{yml1, yml2},
			wantStdout: expectedStylish,
		},
		{
			name:       "compares json files with plain format",
			args:       []string{json1, json2, "--format=plain"},
			wantStdout: expectedPlain,
		},
		{
			name:       "compares yml files with plain format",
			args:       []string{yml1, yml2, "--format=plain"},
			wantStdout: expectedPlain,
		},
		{
			name:       "compares json files with json format",
			args:       []string{json1, json2, "--format=json"},
			wantStdout: expectedJSON,
		},
		{
			name:       "compares yml files with json format",
			args:       []string{yml1, yml2, "--format=json"},
			wantStdout: expectedJSON,
		},
		{
			name:       "accepts short format flag",
			args:       []string{json1, json2, "-f", "plain"},
			wantStdout: expectedPlain,
		},
	}

	for _, tc := range successCases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(t, tc.args...)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantStdout, strings.TrimSpace(stdout))
			assert.Empty(t, stderr)
		})
	}

	errorCases := []struct {
		name       string
		args       []string
		wantStdout string
		wantStderr string
	}{
		{
			name:       "rejects invocation without arguments",
			args:       nil,
			wantStdout: help,
			wantStderr: "incorrect usage: expected 2 arguments, got 0\n",
		},
		{
			name:       "rejects invocation with one argument",
			args:       []string{json1},
			wantStdout: help,
			wantStderr: "incorrect usage: expected 2 arguments, got 1\n",
		},
		{
			name:       "rejects invocation with too many arguments",
			args:       []string{json1, json2, json1},
			wantStdout: help,
			wantStderr: "incorrect usage: expected 2 arguments, got 3\n",
		},
		{
			name:       "rejects unknown output format",
			args:       []string{json1, json2, "--format=xml"},
			wantStdout: help,
			wantStderr: "incorrect usage: invalid format: xml\n",
		},
		{
			name:       "rejects unsupported input extension",
			args:       []string{"a.txt", "b.txt"},
			wantStdout: help,
			wantStderr: "incorrect usage: unsupported extension: '.txt' for file 'a.txt'\n",
		},
		{
			name:       "rejects incompatible input formats",
			args:       []string{json1, yml2},
			wantStdout: help,
			wantStderr: fmt.Sprintf(
				"incorrect usage: files have different formats: files '%s' and '%s'\n",
				json1,
				yml2,
			),
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(t, tc.args...)
			assert.Error(t, err)
			assert.Equal(t, tc.wantStdout, strings.TrimSpace(stdout))
			assert.Equal(t, tc.wantStderr, stderr)
		})
	}
}

func runCommand(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()

	app := NewCommand()
	var outBuf, errBuf bytes.Buffer
	app.Writer = &outBuf
	app.ErrWriter = &errBuf
	app.ExitErrHandler = func(_ context.Context, cmd *urfave.Command, runErr error) {
		if runErr != nil {
			_, writeErr := fmt.Fprintln(cmd.ErrWriter, runErr.Error())
			require.NoError(t, writeErr)
		}
	}

	runArgs := append([]string{"gendiff"}, args...)
	err = app.Run(context.Background(), runArgs)

	return outBuf.String(), errBuf.String(), err
}

func readFixture(t *testing.T, name string) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Clean(filepath.Join(getFixturePath(), name)))
	require.NoError(t, err)

	return strings.TrimSpace(string(content))
}

func getFixturePath() string {
	return filepath.Join("testdata", "fixture")
}
