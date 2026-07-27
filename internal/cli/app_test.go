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
	expected, err := os.ReadFile(filepath.Join(getFixturePath(), "expected.txt"))
	require.NoError(t, err)

	testCases := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "help",
			args: []string{"--help"},
			want: getHelp(t),
		},
		{
			name: "plain json files, stylish format",
			args: []string{filepath.Join(getFixturePath(), "file1.json"), filepath.Join(getFixturePath(), "file2.json")},
			want: string(expected),
		},
		{
			name: "plain yml files, stylish format",
			args: []string{filepath.Join(getFixturePath(), "file1.yml"), filepath.Join(getFixturePath(), "file2.yml")},
			want: string(expected),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := NewCommand()
			app.ExitErrHandler = func(_ context.Context, cmd *urfave.Command, err error) {
				if err != nil {
					_, err := fmt.Fprintln(cmd.ErrWriter, err.Error())
					require.NoError(t, err)
				}
			}

			var stdout, stderr bytes.Buffer
			app.Writer = &stdout
			app.ErrWriter = &stderr

			runArgs := append([]string{"gendiff"}, tc.args...)

			err := app.Run(context.Background(), runArgs)
			assert.NoError(t, err)
			assert.Empty(t, stderr.String())
			assert.Equal(t, tc.want, strings.TrimSpace(stdout.String()))
		})
	}
}

func getHelp(t *testing.T) string {
	t.Helper()
	helpPath := filepath.Join(getFixturePath(), "help.txt")
	helpContent, err := os.ReadFile(filepath.Clean(helpPath))
	require.NoError(t, err)

	return strings.TrimSpace(string(helpContent))
}

func getFixturePath() string {
	return filepath.Join("testdata", "fixture")
}
