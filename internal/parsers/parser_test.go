package parsers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONParser(t *testing.T) {
	p := jsonParser{}

	t.Run("parses nested object with scalar list and null values", func(t *testing.T) {
		got, err := p.parse([]byte(`{
			"foo": "bar",
			"count": 1,
			"ok": true,
			"empty": null,
			"list": ["a", 2],
			"nested": {"key": "value"}
		}`))
		require.NoError(t, err)
		require.Equal(t, map[string]any{
			"foo":   "bar",
			"count": float64(1),
			"ok":    true,
			"empty": nil,
			"list":  []any{"a", float64(2)},
			"nested": map[string]any{
				"key": "value",
			},
		}, got)
	})

	t.Run("returns error for malformed content", func(t *testing.T) {
		got, err := p.parse([]byte(`{not-json`))
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("returns error when root value is not an object", func(t *testing.T) {
		got, err := p.parse([]byte(`[1, 2]`))
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestYAMLParser(t *testing.T) {
	p := yamlParser{}

	t.Run("parses nested object with scalar list and null values", func(t *testing.T) {
		got, err := p.parse([]byte(`
foo: bar
count: 1
ok: true
empty: null
list:
  - a
  - 2
nested:
  key: value
`))
		require.NoError(t, err)
		require.Equal(t, map[string]any{
			"foo":   "bar",
			"count": 1,
			"ok":    true,
			"empty": nil,
			"list":  []any{"a", 2},
			"nested": map[string]any{
				"key": "value",
			},
		}, got)
	})

	t.Run("returns error for malformed content", func(t *testing.T) {
		got, err := p.parse([]byte(":\n  - broken"))
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestParseFiles(t *testing.T) {
	dir := t.TempDir()

	json1 := writeTempFile(t, dir, "a.json", `{"foo":"bar","n":1}`)
	json2 := writeTempFile(t, dir, "b.json", `{"foo":"baz","n":2}`)
	yml1 := writeTempFile(t, dir, "a.yml", "foo: bar\nn: 1\n")
	yml2 := writeTempFile(t, dir, "b.yml", "foo: baz\nn: 2\n")
	yaml1 := writeTempFile(t, dir, "a.yaml", "foo: bar\nn: 1\n")
	yaml2 := writeTempFile(t, dir, "b.yaml", "foo: baz\nn: 2\n")
	jsonUpper := writeTempFile(t, dir, "upper1.JSON", `{"foo":"bar"}`)
	jsonUpper2 := writeTempFile(t, dir, "upper2.JSON", `{"foo":"baz"}`)

	testCases := []struct {
		name  string
		path1 string
		path2 string
		want1 map[string]any
		want2 map[string]any
	}{
		{
			name:  "parses pair of json configs",
			path1: json1,
			path2: json2,
			want1: map[string]any{"foo": "bar", "n": float64(1)},
			want2: map[string]any{"foo": "baz", "n": float64(2)},
		},
		{
			name:  "parses pair of yml configs",
			path1: yml1,
			path2: yml2,
			want1: map[string]any{"foo": "bar", "n": 1},
			want2: map[string]any{"foo": "baz", "n": 2},
		},
		{
			name:  "parses pair of yaml configs",
			path1: yaml1,
			path2: yaml2,
			want1: map[string]any{"foo": "bar", "n": 1},
			want2: map[string]any{"foo": "baz", "n": 2},
		},
		{
			name:  "allows mixing yml and yaml extensions",
			path1: yml1,
			path2: yaml2,
			want1: map[string]any{"foo": "bar", "n": 1},
			want2: map[string]any{"foo": "baz", "n": 2},
		},
		{
			name:  "treats extension case-insensitively",
			path1: jsonUpper,
			path2: jsonUpper2,
			want1: map[string]any{"foo": "bar"},
			want2: map[string]any{"foo": "baz"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got1, got2, err := ParseFiles(tc.path1, tc.path2)
			assert.NoError(t, err)
			assert.Equal(t, tc.want1, got1)
			assert.Equal(t, tc.want2, got2)
		})
	}
}

func TestParseFiles_Errors(t *testing.T) {
	dir := t.TempDir()

	validJSON := writeTempFile(t, dir, "valid.json", `{"ok":true}`)
	validYML := writeTempFile(t, dir, "valid.yml", "ok: true\n")
	invalidJSON := writeTempFile(t, dir, "invalid.json", `{broken`)
	invalidYML := writeTempFile(t, dir, "invalid.yml", ":\n  - broken\n")
	emptyJSON := writeTempFile(t, dir, "empty.json", "")
	noExt := writeTempFile(t, dir, "noext", `{"ok":true}`)
	unsupported := writeTempFile(t, dir, "file.txt", "hello")
	subdir := filepath.Join(dir, "subdir.json")
	require.NoError(t, os.Mkdir(subdir, 0o700))

	unreadable := writeTempFile(t, dir, "unreadable.json", `{"ok":true}`)
	require.NoError(t, os.Chmod(unreadable, 0o000))
	t.Cleanup(func() {
		require.NoError(t, os.Chmod(unreadable, 0o600))
	})

	testCases := []struct {
		name     string
		path1    string
		path2    string
		wantErrs []error
	}{
		{
			name:     "rejects path without extension",
			path1:    noExt,
			path2:    validJSON,
			wantErrs: []error{ErrPathHasNoExtension},
		},
		{
			name:     "rejects unsupported extension",
			path1:    unsupported,
			path2:    validJSON,
			wantErrs: []error{ErrUnsupportedExtension},
		},
		{
			name:     "rejects incompatible formats",
			path1:    validJSON,
			path2:    validYML,
			wantErrs: []error{ErrDifferentFileFormats},
		},
		{
			name:     "rejects empty file",
			path1:    emptyJSON,
			path2:    validJSON,
			wantErrs: []error{ErrFileIsEmpty},
		},
		{
			name:     "rejects directory path",
			path1:    subdir,
			path2:    validJSON,
			wantErrs: []error{ErrNotRegularFile},
		},
		{
			name:     "rejects unparsable json content",
			path1:    invalidJSON,
			path2:    validJSON,
			wantErrs: []error{ErrFailedToParseFile},
		},
		{
			name:     "rejects unparsable yaml content",
			path1:    invalidYML,
			path2:    validYML,
			wantErrs: []error{ErrFailedToParseFile},
		},
		{
			name:     "rejects unreadable file",
			path1:    unreadable,
			path2:    validJSON,
			wantErrs: []error{ErrFailedToReadFile, os.ErrPermission},
		},
		{
			name:     "rejects missing file",
			path1:    filepath.Join(dir, "missing.json"),
			path2:    validJSON,
			wantErrs: []error{os.ErrNotExist},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got1, got2, err := ParseFiles(tc.path1, tc.path2)
			assert.Error(t, err)
			assert.Nil(t, got1)
			assert.Nil(t, got2)

			for _, wantErr := range tc.wantErrs {
				assert.ErrorIs(t, err, wantErr)
			}
		})
	}
}

func TestExtractExt(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		want    string
		wantErr error
	}{
		{
			name: "returns lowercase json extension",
			path: "file.json",
			want: ".json",
		},
		{
			name: "normalizes uppercase json extension",
			path: "file.JSON",
			want: ".json",
		},
		{
			name: "returns yml extension",
			path: "file.yml",
			want: ".yml",
		},
		{
			name: "returns yaml extension",
			path: "file.yaml",
			want: ".yaml",
		},
		{
			name: "normalizes uppercase extension in nested path",
			path: "dir/file.YML",
			want: ".yml",
		},
		{
			name:    "fails when extension is missing",
			path:    "noext",
			wantErr: ErrPathHasNoExtension,
		},
		{
			name:    "fails when nested path has no extension",
			path:    "dir/noext",
			wantErr: ErrPathHasNoExtension,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := extractExt(tc.path)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Empty(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFileFormat(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		want    format
		wantErr error
	}{
		{
			name: "detects json format",
			path: "a.json",
			want: formatJSON,
		},
		{
			name: "detects json format from uppercase extension",
			path: "a.JSON",
			want: formatJSON,
		},
		{
			name: "detects yaml format from yml extension",
			path: "a.yml",
			want: formatYAML,
		},
		{
			name: "detects yaml format from yaml extension",
			path: "a.yaml",
			want: formatYAML,
		},
		{
			name:    "rejects unsupported extension",
			path:    "a.txt",
			wantErr: ErrUnsupportedExtension,
		},
		{
			name:    "rejects path without extension",
			path:    "a",
			wantErr: ErrPathHasNoExtension,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fileFormat(tc.path)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Empty(t, got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	return path
}
