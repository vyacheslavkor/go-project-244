package parsers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type parser interface {
	parse(content []byte) (any, error)
}

type format string

const (
	formatJSON format = "json"
	formatYAML format = "yaml"
)

var extFormats = map[string]format{
	".json": formatJSON,
	".yml":  formatYAML,
	".yaml": formatYAML,
}

var formatParsers = map[format]parser{
	formatJSON: jsonParser{},
	formatYAML: yamlParser{},
}

var (
	// ErrMissingExtension is returned when a path has no file extension.
	ErrMissingExtension = errors.New("file path has no extension")
	// ErrUnsupportedExtension is returned when the file extension is not supported.
	ErrUnsupportedExtension = errors.New("unsupported extension")
	// ErrNotRegularFile is returned when the path is not a regular file.
	ErrNotRegularFile = errors.New("path is not a regular file")
	// ErrEmptyFile is returned when the file has zero size.
	ErrEmptyFile = errors.New("file is empty")
	// ErrInvalidRoot is returned when the root value is not a JSON object or YAML mapping.
	ErrInvalidRoot = errors.New("root value must be a JSON object or a YAML mapping")
	// ErrParseFile is returned when file contents cannot be decoded.
	ErrParseFile = errors.New("failed to parse file")
	// ErrReadFile is returned when the file cannot be read from disk.
	ErrReadFile = errors.New("failed to read file")
)

// ParseFiles reads and parses two configuration files and returns their
// contents as nested maps.
//
// beforePath and afterPath must point to existing non-empty regular files
// whose root value is a JSON object or YAML mapping. Supported formats are JSON
// (.json) and YAML (.yml, .yaml). The files can use different formats.
//
// Paths are validated as regular non-empty files before extension checks,
// so directories and missing paths are reported as file errors rather than
// "no extension" usage errors.
//
// After decoding, values are normalized for comparison: integers to float64,
// time.Time to a string (date-only as YYYY-MM-DD, otherwise RFC3339).
//
// On failure, err may wrap one of the package sentinel errors
// (for example [ErrUnsupportedExtension], [ErrInvalidRoot], [ErrParseFile]).
func ParseFiles(beforePath, afterPath string) (before, after map[string]any, err error) {
	if err := validateInputFile(beforePath); err != nil {
		return nil, nil, err
	}

	if err := validateInputFile(afterPath); err != nil {
		return nil, nil, err
	}

	beforeParser, err := resolveParser(beforePath)
	if err != nil {
		return nil, nil, err
	}

	afterParser, err := resolveParser(afterPath)
	if err != nil {
		return nil, nil, err
	}

	beforeContent, err := readFileContent(beforePath)
	if err != nil {
		return nil, nil, err
	}

	afterContent, err := readFileContent(afterPath)
	if err != nil {
		return nil, nil, err
	}

	before, err = parseFile(beforeParser, beforePath, beforeContent)
	if err != nil {
		return nil, nil, err
	}

	after, err = parseFile(afterParser, afterPath, afterContent)
	if err != nil {
		return nil, nil, err
	}

	return before, after, nil
}

func parseFile(p parser, path string, content []byte) (map[string]any, error) {
	parsed, err := p.parse(content)
	if err != nil {
		return nil, fmt.Errorf("%w: '%s': %w", ErrParseFile, path, err)
	}

	root, err := rootMap(parsed)
	if err != nil {
		return nil, fmt.Errorf("file '%s': %w", path, err)
	}

	normalizeMap(root)

	return root, nil
}

func rootMap(root any) (map[string]any, error) {
	if root == nil {
		return nil, fmt.Errorf("%w: got null", ErrInvalidRoot)
	}

	m, ok := root.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%w: got %s", ErrInvalidRoot, rootKind(root))
	}

	return m, nil
}

func rootKind(root any) string {
	switch root.(type) {
	case []any:
		return "array"
	case string:
		return "string"
	case bool:
		return "boolean"
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return "number"
	default:
		return fmt.Sprintf("%T", root)
	}
}

func resolveParser(filePath string) (parser, error) {
	format, err := fileFormat(filePath)
	if err != nil {
		return nil, err
	}

	return formatParsers[format], nil
}

func fileFormat(path string) (format, error) {
	ext, err := extractExt(path)
	if err != nil {
		return "", err
	}

	f, ok := extFormats[ext]
	if !ok {
		return "", fmt.Errorf("%w: '%s' for file '%s'", ErrUnsupportedExtension, ext, path)
	}

	return f, nil
}

func extractExt(path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		return "", fmt.Errorf("%w: '%s'", ErrMissingExtension, path)
	}

	return ext, nil
}

func validateInputFile(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%w: '%s': %w", ErrReadFile, path, err)
	}

	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("%w: '%s'", ErrNotRegularFile, path)
	}

	if fileInfo.Size() == 0 {
		return fmt.Errorf("%w: '%s'", ErrEmptyFile, path)
	}

	return nil
}

func readFileContent(path string) ([]byte, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("%w: '%s': %w", ErrReadFile, path, err)
	}

	return content, nil
}
