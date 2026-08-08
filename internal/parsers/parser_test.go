package parsers

import (
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeParser struct {
	value any
}

func (p fakeParser) parse(content []byte) (any, error) {
	return p.value, nil
}

func TestJSONParser(t *testing.T) {
	p := jsonParser{}

	t.Run("parses nested object with scalar list and null values", func(t *testing.T) {
		got, err := parseFile(p, "file.json", []byte(`{
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
		got, err := parseFile(p, "file.json", []byte(`{not-json`))
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("returns error when root value is not an object", func(t *testing.T) {
		got, err := parseFile(p, "file.json", []byte(`[1, 2]`))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidRoot)
		require.Nil(t, got)
	})

	t.Run("returns error when root value is null", func(t *testing.T) {
		got, err := parseFile(p, "file.json", []byte(`null`))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidRoot)
		require.Nil(t, got)
	})

	t.Run("returns error when root value is a scalar", func(t *testing.T) {
		got, err := parseFile(p, "file.json", []byte(`"hello"`))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidRoot)
		require.Nil(t, got)
	})
}

func TestYAMLParser(t *testing.T) {
	p := yamlParser{}

	t.Run("parses nested object with scalar list and null values", func(t *testing.T) {
		got, err := parseFile(p, "file.yml", []byte(`
foo: bar
count: 1
ok: true
empty: null
list:
  - a
  - 1
  - 1.0
nested:
  key: value
`))
		require.NoError(t, err)
		require.Equal(t, map[string]any{
			"foo":   "bar",
			"count": float64(1),
			"ok":    true,
			"empty": nil,
			"list":  []any{"a", float64(1), float64(1)},
			"nested": map[string]any{
				"key": "value",
			},
		}, got)
	})

	t.Run("returns error for malformed content", func(t *testing.T) {
		got, err := parseFile(p, "file.yml", []byte(":\n  - broken"))
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("returns error when root value is a sequence", func(t *testing.T) {
		got, err := parseFile(p, "file.yml", []byte("- a\n- b\n"))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidRoot)
		require.Nil(t, got)
	})

	t.Run("returns error when root value is null", func(t *testing.T) {
		got, err := parseFile(p, "file.yml", []byte("null\n"))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidRoot)
		require.Nil(t, got)
	})

	t.Run("returns error when root value is a scalar", func(t *testing.T) {
		got, err := parseFile(p, "file.yml", []byte("42\n"))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidRoot)
		require.Nil(t, got)
	})

	t.Run("the time is parsed as a string", func(t *testing.T) {
		got, err := parseFile(p, "file.yml", []byte("released: 2023-01-01\n"))
		require.NoError(t, err)
		require.Equal(t, map[string]any{
			"released": "2023-01-01",
		}, got)
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
			want1: map[string]any{"foo": "bar", "n": float64(1)},
			want2: map[string]any{"foo": "baz", "n": float64(2)},
		},
		{
			name:  "parses pair of yaml configs",
			path1: yaml1,
			path2: yaml2,
			want1: map[string]any{"foo": "bar", "n": float64(1)},
			want2: map[string]any{"foo": "baz", "n": float64(2)},
		},
		{
			name:  "allows mixing yml and yaml extensions",
			path1: yml1,
			path2: yaml2,
			want1: map[string]any{"foo": "bar", "n": float64(1)},
			want2: map[string]any{"foo": "baz", "n": float64(2)},
		},
		{
			name:  "treats extension case-insensitively",
			path1: jsonUpper,
			path2: jsonUpper2,
			want1: map[string]any{"foo": "bar"},
			want2: map[string]any{"foo": "baz"},
		},
		{
			name:  "parses pair of json and yml configs",
			path1: json1,
			path2: yml2,
			want1: map[string]any{"foo": "bar", "n": float64(1)},
			want2: map[string]any{"foo": "baz", "n": float64(2)},
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
	subdirNoExt := filepath.Join(dir, "subdir")
	require.NoError(t, os.Mkdir(subdirNoExt, 0o700))

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
			wantErrs: []error{ErrMissingExtension},
		},
		{
			name:     "rejects unsupported extension",
			path1:    unsupported,
			path2:    validJSON,
			wantErrs: []error{ErrUnsupportedExtension},
		},
		{
			name:     "rejects empty file",
			path1:    emptyJSON,
			path2:    validJSON,
			wantErrs: []error{ErrEmptyFile},
		},
		{
			name:     "rejects directory path",
			path1:    subdir,
			path2:    validJSON,
			wantErrs: []error{ErrNotRegularFile},
		},
		{
			name:     "rejects directory path without extension as not a regular file",
			path1:    subdirNoExt,
			path2:    validJSON,
			wantErrs: []error{ErrNotRegularFile},
		},
		{
			name:     "rejects unparsable json content",
			path1:    invalidJSON,
			path2:    validJSON,
			wantErrs: []error{ErrParseFile},
		},
		{
			name:     "rejects unparsable yaml content",
			path1:    invalidYML,
			path2:    validYML,
			wantErrs: []error{ErrParseFile},
		},
		{
			name:     "rejects json root array",
			path1:    writeTempFile(t, dir, "array.json", `[1,2]`),
			path2:    validJSON,
			wantErrs: []error{ErrInvalidRoot},
		},
		{
			name:     "rejects yaml root sequence",
			path1:    writeTempFile(t, dir, "array.yml", "- a\n- b\n"),
			path2:    validYML,
			wantErrs: []error{ErrInvalidRoot},
		},
		{
			name:     "rejects unreadable file",
			path1:    unreadable,
			path2:    validJSON,
			wantErrs: []error{ErrReadFile, os.ErrPermission},
		},
		{
			name:     "rejects missing file",
			path1:    filepath.Join(dir, "missing.json"),
			path2:    validJSON,
			wantErrs: []error{ErrReadFile, os.ErrNotExist},
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
			wantErr: ErrMissingExtension,
		},
		{
			name:    "fails when nested path has no extension",
			path:    "dir/noext",
			wantErr: ErrMissingExtension,
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

func TestDetectFormat(t *testing.T) {
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
			wantErr: ErrMissingExtension,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := detectFormat(tc.path)
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

func TestParseFile_normalizesValues(t *testing.T) {
	val := map[string]any{
		"uint":     uint64(1),
		"uint_max": uint64(math.MaxUint64),
		"nested": map[string]any{
			"int":      int(2),
			"released": time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			"updated":  time.Date(2023, 1, 1, 15, 4, 5, 0, time.UTC),
		},
		"list": []any{
			uint64(math.MaxUint64),
			time.Date(2023, 1, 1, 15, 4, 5, 0, time.UTC),
		},
	}

	p := fakeParser{value: val}

	got, err := parseFile(p, "fake", []byte(`{"fake": "fake"}`))
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"uint":     float64(1),
		"uint_max": float64(math.MaxUint64),
		"nested": map[string]any{
			"int":      float64(2),
			"released": "2023-01-01",
			"updated":  "2023-01-01T15:04:05Z",
		},
		"list": []any{
			float64(math.MaxUint64),
			"2023-01-01T15:04:05Z",
		},
	}, got)
}

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	return path
}
