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

type commandFixtures struct {
	json1, json2    string
	yml1, yml2      string
	txt1, txt2      string
	emptyJSON       string
	arrayJSON       string
	arrayYML        string
	dirAsJSON       string
	badJSON         string
	expectedStylish string
	expectedPlain   string
	expectedJSON    string
	help            string
}

func setupCommandFixtures(t *testing.T) commandFixtures {
	t.Helper()

	fixtureDir := getFixturePath()
	txtDir := t.TempDir()

	txt1 := filepath.Join(txtDir, "a.txt")
	txt2 := filepath.Join(txtDir, "b.txt")
	require.NoError(t, os.WriteFile(txt1, []byte("hello"), 0o600))
	require.NoError(t, os.WriteFile(txt2, []byte("world"), 0o600))

	emptyJSON := filepath.Join(txtDir, "empty.json")
	require.NoError(t, os.WriteFile(emptyJSON, []byte{}, 0o600))

	arrayJSON := filepath.Join(txtDir, "array.json")
	require.NoError(t, os.WriteFile(arrayJSON, []byte(`[1,2]`), 0o600))

	arrayYML := filepath.Join(txtDir, "array.yml")
	require.NoError(t, os.WriteFile(arrayYML, []byte("- a\n- b\n"), 0o600))

	dirAsJSON := filepath.Join(txtDir, "subdir.json")
	require.NoError(t, os.Mkdir(dirAsJSON, 0o700))

	badJSON := filepath.Join(txtDir, "bad.json")
	require.NoError(t, os.WriteFile(badJSON, []byte(`{broken`), 0o600))

	return commandFixtures{
		json1:           filepath.Join(fixtureDir, "file1.json"),
		json2:           filepath.Join(fixtureDir, "file2.json"),
		yml1:            filepath.Join(fixtureDir, "file1.yml"),
		yml2:            filepath.Join(fixtureDir, "file2.yml"),
		txt1:            txt1,
		txt2:            txt2,
		emptyJSON:       emptyJSON,
		arrayJSON:       arrayJSON,
		arrayYML:        arrayYML,
		dirAsJSON:       dirAsJSON,
		badJSON:         badJSON,
		expectedStylish: readFixture(t, "expected_stylish.txt"),
		expectedPlain:   readFixture(t, "expected_plain.txt"),
		expectedJSON:    readFixture(t, "expected_json.txt"),
		help:            readFixture(t, "help.txt"),
	}
}

func TestNewCommandSuccess(t *testing.T) {
	f := setupCommandFixtures(t)

	cases := []struct {
		name       string
		args       []string
		wantStdout string
	}{
		{
			name:       "prints help when help flag is passed",
			args:       []string{"--help"},
			wantStdout: f.help,
		},
		{
			name:       "compares json files with default stylish format",
			args:       []string{f.json1, f.json2},
			wantStdout: f.expectedStylish,
		},
		{
			name:       "compares yml files with default stylish format",
			args:       []string{f.yml1, f.yml2},
			wantStdout: f.expectedStylish,
		},
		{
			name:       "compares json files with plain format",
			args:       []string{f.json1, f.json2, "--format=plain"},
			wantStdout: f.expectedPlain,
		},
		{
			name:       "compares yml files with plain format",
			args:       []string{f.yml1, f.yml2, "--format=plain"},
			wantStdout: f.expectedPlain,
		},
		{
			name:       "compares json files with json format",
			args:       []string{f.json1, f.json2, "--format=json"},
			wantStdout: f.expectedJSON,
		},
		{
			name:       "compares yml files with json format",
			args:       []string{f.yml1, f.yml2, "--format=json"},
			wantStdout: f.expectedJSON,
		},
		{
			name:       "accepts short format flag",
			args:       []string{f.json1, f.json2, "-f", "plain"},
			wantStdout: f.expectedPlain,
		},
		{
			name:       "plain format prints nothing for identical nested files",
			args:       []string{f.json1, f.json1, "--format=plain"},
			wantStdout: "",
		},
		{
			name:       "comparing files with different extensions",
			args:       []string{f.json1, f.yml2},
			wantStdout: f.expectedStylish,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(t, tc.args...)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantStdout, strings.TrimSpace(stdout))
			assert.Empty(t, stderr)
		})
	}
}

func TestNewCommandUsageErrors(t *testing.T) {
	f := setupCommandFixtures(t)

	cases := []struct {
		name       string
		args       []string
		wantStdout string
		wantStderr string
	}{
		{
			name:       "rejects invocation without arguments",
			args:       nil,
			wantStdout: f.help,
			wantStderr: "incorrect usage: expected 2 arguments, got 0\n",
		},
		{
			name:       "rejects invocation with one argument",
			args:       []string{f.json1},
			wantStdout: f.help,
			wantStderr: "incorrect usage: expected 2 arguments, got 1\n",
		},
		{
			name:       "rejects invocation with too many arguments",
			args:       []string{f.json1, f.json2, f.json1},
			wantStdout: f.help,
			wantStderr: "incorrect usage: expected 2 arguments, got 3\n",
		},
		{
			name:       "rejects unknown output format",
			args:       []string{f.json1, f.json2, "--format=xml"},
			wantStdout: f.help,
			wantStderr: "incorrect usage: invalid format: xml\n",
		},
		{
			name:       "rejects unsupported input extension",
			args:       []string{f.txt1, f.txt2},
			wantStdout: f.help,
			wantStderr: fmt.Sprintf("incorrect usage: unsupported extension: '.txt' for file '%s'\n", f.txt1),
		},

		{
			name:       "rejects empty file as usage error",
			args:       []string{f.emptyJSON, f.json2},
			wantStdout: f.help,
			wantStderr: fmt.Sprintf("incorrect usage: file is empty: '%s'\n", f.emptyJSON),
		},
		{
			name:       "rejects directory path as usage error",
			args:       []string{f.dirAsJSON, f.json2},
			wantStdout: f.help,
			wantStderr: fmt.Sprintf("incorrect usage: path is not a regular file: '%s'\n", f.dirAsJSON),
		},
		{
			name:       "rejects json root array as usage error",
			args:       []string{f.arrayJSON, f.json2},
			wantStdout: f.help,
			wantStderr: fmt.Sprintf(
				"incorrect usage: file '%s': root value must be a JSON object or a YAML mapping: got array\n",
				f.arrayJSON,
			),
		},
		{
			name:       "rejects yaml root sequence as usage error",
			args:       []string{f.arrayYML, f.yml2},
			wantStdout: f.help,
			wantStderr: fmt.Sprintf(
				"incorrect usage: file '%s': root value must be a JSON object or a YAML mapping: got array\n",
				f.arrayYML,
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(t, tc.args...)
			assert.Error(t, err)
			assert.Equal(t, tc.wantStdout, strings.TrimSpace(stdout))
			assert.Equal(t, tc.wantStderr, stderr)
		})
	}
}

func TestNewCommandOperationalErrors(t *testing.T) {
	f := setupCommandFixtures(t)

	cases := []struct {
		name       string
		args       []string
		wantStderr string
	}{
		{
			name:       "rejects missing file without dumping help",
			args:       []string{"missing1.json", "missing2.json"},
			wantStderr: "failed to read file: 'missing1.json': stat missing1.json: no such file or directory\n",
		},
		{
			name: "rejects malformed json without dumping help",
			args: []string{f.badJSON, f.json2},
			wantStderr: fmt.Sprintf(
				"failed to parse file: '%s': invalid character 'b' looking for beginning of object key string\n",
				f.badJSON,
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(t, tc.args...)
			assert.Error(t, err)
			assert.Empty(t, stdout)
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
