package parsers

import (
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

	t.Run("returns error when root value is a boolean", func(t *testing.T) {
		got, err := parseFile(p, "file.json", []byte(`true`))
		require.ErrorIs(t, err, ErrInvalidRoot)
		require.ErrorContains(t, err, "got boolean")
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

	t.Run("returns error when root value is a timestamp", func(t *testing.T) {
		got, err := parseFile(p, "file.yml", []byte("2023-01-01\n"))
		require.ErrorIs(t, err, ErrInvalidRoot)
		require.ErrorContains(t, err, "got timestamp")
		require.Nil(t, got)
	})

	t.Run("returns error when root mapping has non-string keys", func(t *testing.T) {
		got, err := parseFile(p, "file.yml", []byte("1: a\n"))
		require.ErrorIs(t, err, ErrInvalidRoot)
		require.ErrorContains(t, err, "got mapping with non-string keys")
		require.Nil(t, got)
	})

	t.Run("normalizes a timestamp to a string", func(t *testing.T) {
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
		{
			name:  "parses pair of yml and json configs",
			path1: yml1,
			path2: json2,
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

func TestParseFilesErrors(t *testing.T) {
	dir := t.TempDir()

	validJSON := writeTempFile(t, dir, "valid.json", `{"ok":true}`)
	validYML := writeTempFile(t, dir, "valid.yml", "ok: true\n")
	invalidJSON := writeTempFile(t, dir, "invalid.json", `{broken`)
	invalidYML := writeTempFile(t, dir, "invalid.yml", ":\n  - broken\n")
	emptyJSON := writeTempFile(t, dir, "empty.json", "")
	noExt := writeTempFile(t, dir, "noext", `{"ok":true}`)
	unsupported := writeTempFile(t, dir, "file.txt", "hello")
	arrayJSON := writeTempFile(t, dir, "array.json", `[1,2]`)
	arrayYML := writeTempFile(t, dir, "array.yml", "- a\n- b\n")
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
			name:     "rejects first path without extension",
			path1:    noExt,
			path2:    validJSON,
			wantErrs: []error{ErrMissingExtension},
		},
		{
			name:     "rejects second path without extension",
			path1:    validJSON,
			path2:    noExt,
			wantErrs: []error{ErrMissingExtension},
		},
		{
			name:     "rejects unsupported extension on first file",
			path1:    unsupported,
			path2:    validJSON,
			wantErrs: []error{ErrUnsupportedExtension},
		},
		{
			name:     "rejects unsupported extension on second file",
			path1:    validJSON,
			path2:    unsupported,
			wantErrs: []error{ErrUnsupportedExtension},
		},
		{
			name:     "rejects empty first file",
			path1:    emptyJSON,
			path2:    validJSON,
			wantErrs: []error{ErrEmptyFile},
		},
		{
			name:     "rejects empty second file",
			path1:    validJSON,
			path2:    emptyJSON,
			wantErrs: []error{ErrEmptyFile},
		},
		{
			name:     "rejects directory in first path",
			path1:    subdir,
			path2:    validJSON,
			wantErrs: []error{ErrNotRegularFile},
		},
		{
			name:     "rejects directory in second path",
			path1:    validJSON,
			path2:    subdir,
			wantErrs: []error{ErrNotRegularFile},
		},
		{
			name:     "rejects directory without extension in first path as not a regular file",
			path1:    subdirNoExt,
			path2:    validJSON,
			wantErrs: []error{ErrNotRegularFile},
		},
		{
			name:     "rejects directory without extension in second path as not a regular file",
			path1:    validJSON,
			path2:    subdirNoExt,
			wantErrs: []error{ErrNotRegularFile},
		},
		{
			name:     "rejects unparsable json content in first file",
			path1:    invalidJSON,
			path2:    validJSON,
			wantErrs: []error{ErrParseFile},
		},
		{
			name:     "rejects unparsable json content in second file",
			path1:    validJSON,
			path2:    invalidJSON,
			wantErrs: []error{ErrParseFile},
		},
		{
			name:     "rejects unparsable yaml content in first file",
			path1:    invalidYML,
			path2:    validYML,
			wantErrs: []error{ErrParseFile},
		},
		{
			name:     "rejects unparsable yaml content in second file",
			path1:    validYML,
			path2:    invalidYML,
			wantErrs: []error{ErrParseFile},
		},
		{
			name:     "rejects json root array in first file",
			path1:    arrayJSON,
			path2:    validJSON,
			wantErrs: []error{ErrInvalidRoot},
		},
		{
			name:     "rejects json root array in second file",
			path1:    validJSON,
			path2:    arrayJSON,
			wantErrs: []error{ErrInvalidRoot},
		},
		{
			name:     "rejects yaml root sequence in first file",
			path1:    arrayYML,
			path2:    validYML,
			wantErrs: []error{ErrInvalidRoot},
		},
		{
			name:     "rejects yaml root sequence in second file",
			path1:    validYML,
			path2:    arrayYML,
			wantErrs: []error{ErrInvalidRoot},
		},
		{
			name:     "rejects unreadable first file",
			path1:    unreadable,
			path2:    validJSON,
			wantErrs: []error{ErrReadFile, os.ErrPermission},
		},
		{
			name:     "rejects unreadable second file",
			path1:    validJSON,
			path2:    unreadable,
			wantErrs: []error{ErrReadFile, os.ErrPermission},
		},
		{
			name:     "rejects missing first file",
			path1:    filepath.Join(dir, "missing.json"),
			path2:    validJSON,
			wantErrs: []error{ErrReadFile, os.ErrNotExist},
		},
		{
			name:     "rejects missing second file",
			path1:    validJSON,
			path2:    filepath.Join(dir, "missing.json"),
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

func TestDescribeRootKind(t *testing.T) {
	testCases := []struct {
		name string
		root any
		want string
	}{
		{
			name: "describes a sequence",
			root: []any{1, 2},
			want: "array",
		},
		{
			name: "describes a string",
			root: "hello",
			want: "string",
		},
		{
			name: "describes a boolean",
			root: true,
			want: "boolean",
		},
		{
			name: "describes an integer",
			root: 42,
			want: "number",
		},
		{
			name: "describes a float",
			root: 4.2,
			want: "number",
		},
		{
			name: "describes a timestamp",
			root: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want: "timestamp",
		},
		{
			name: "describes a mapping with non-string keys",
			root: map[any]any{1: "a"},
			want: "mapping with non-string keys",
		},
		{
			name: "never leaks a go type for a value no parser produces",
			root: struct{}{},
			want: "unsupported value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, describeRootKind(tc.root))
		})
	}
}

func TestParseFileNormalizesValues(t *testing.T) {
	t.Run("converts integers and timestamps at every depth", func(t *testing.T) {
		p := fakeParser{value: map[string]any{
			"uint":     uint64(1),
			"uint_max": uint64(18446744073709551615),
			"nested": map[string]any{
				"int":      int(2),
				"released": time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				"updated":  time.Date(2023, 1, 1, 15, 4, 5, 0, time.UTC),
			},
			"list": []any{
				uint64(18446744073709551615),
				time.Date(2023, 1, 1, 15, 4, 5, 0, time.UTC),
			},
		}}

		got, err := parseFile(p, "fake", []byte(`{"fake": "fake"}`))
		require.NoError(t, err)
		require.Equal(t, map[string]any{
			"uint":     float64(1),
			"uint_max": float64(18446744073709551616),
			"nested": map[string]any{
				"int":      float64(2),
				"released": "2023-01-01",
				"updated":  "2023-01-01 15:04:05",
			},
			"list": []any{
				float64(18446744073709551616),
				"2023-01-01 15:04:05",
			},
		}, got)
	})

	t.Run("keeps every part of a timestamp that tells two values apart", func(t *testing.T) {
		testCases := []struct {
			name string
			val  time.Time
			want string
		}{
			{
				name: "keeps fractional seconds",
				val:  time.Date(2023, 1, 1, 10, 20, 30, 500_000_000, time.UTC),
				want: "2023-01-01T10:20:30.5Z",
			},
			{
				name: "tells a different fraction of a second apart",
				val:  time.Date(2023, 1, 1, 10, 20, 30, 900_000_000, time.UTC),
				want: "2023-01-01T10:20:30.9Z",
			},
			{
				name: "keeps a zone offset at midnight",
				val:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.FixedZone("", 3*60*60)),
				want: "2023-01-01T00:00:00+03:00",
			},
			{
				name: "keeps a zone offset together with fractional seconds",
				val:  time.Date(2023, 1, 1, 10, 20, 30, 500_000_000, time.FixedZone("", 3*60*60)),
				want: "2023-01-01T10:20:30.5+03:00",
			},
			{
				name: "shortens midnight in utc to a date, the way a json string spells it",
				val:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				want: "2023-01-01",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				p := fakeParser{value: map[string]any{"ts": tc.val}}

				got, err := parseFile(p, "fake", nil)
				require.NoError(t, err)
				assert.Equal(t, map[string]any{"ts": tc.want}, got)
			})
		}
	})
}

func TestParseFileNormalizesDateLikeStrings(t *testing.T) {
	testCases := []struct {
		name string
		val  string
		want string
	}{
		{
			name: "reads a date with a two digit month and day",
			val:  "2023-01-01",
			want: "2023-01-01",
		},
		{
			name: "reads a date with a one digit month and day",
			val:  "2023-1-1",
			want: "2023-01-01",
		},
		{
			name: "shortens midnight written as a date and time",
			val:  "2023-01-01 00:00:00",
			want: "2023-01-01",
		},
		{
			name: "shortens midnight written with one digit parts",
			val:  "2023-1-1 0:0:0",
			want: "2023-01-01",
		},
		{
			name: "shortens midnight written in rfc 3339",
			val:  "2023-01-01T00:00:00Z",
			want: "2023-01-01",
		},
		{
			name: "reads a date and time with two digit parts",
			val:  "2023-01-01 15:04:05",
			want: "2023-01-01 15:04:05",
		},
		{
			name: "reads a date and time with one digit parts",
			val:  "2023-1-1 15:4:5",
			want: "2023-01-01 15:04:05",
		},
		{
			name: "keeps fractional seconds",
			val:  "2023-01-01T10:20:30.5Z",
			want: "2023-01-01T10:20:30.5Z",
		},
		{
			name: "tells a different fraction of a second apart",
			val:  "2023-01-01T10:20:30.9Z",
			want: "2023-01-01T10:20:30.9Z",
		},
		{
			name: "keeps a zone offset at midnight",
			val:  "2023-01-01T00:00:00+03:00",
			want: "2023-01-01T00:00:00+03:00",
		},
		{
			name: "keeps a zone offset together with fractional seconds",
			val:  "2023-01-01T10:20:30.5+03:00",
			want: "2023-01-01T10:20:30.5+03:00",
		},
		{
			name: "leaves a string that is not a date",
			val:  "hello",
			want: "hello",
		},
		{
			name: "leaves a date with impossible month and day",
			val:  "2023-13-45",
			want: "2023-13-45",
		},
		{
			name: "leaves a time without a date",
			val:  "15:04:05",
			want: "15:04:05",
		},
		{
			name: "leaves a lowercase rfc 3339 separator, as yaml does",
			val:  "2023-01-01t15:04:05z",
			want: "2023-01-01t15:04:05z",
		},
		{
			name: "leaves a date and time with a separate offset, as yaml does",
			val:  "2023-01-01 15:04:05 +03:00",
			want: "2023-01-01 15:04:05 +03:00",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := fakeParser{value: map[string]any{"ts": tc.val}}

			got, err := parseFile(p, "fake", nil)
			require.NoError(t, err)
			assert.Equal(t, map[string]any{"ts": tc.want}, got)
		})
	}
}

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	return path
}
