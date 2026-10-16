package parsers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type parser interface {
	parse(content []byte) (any, error)
}

var extParsers = map[string]parser{
	".json": jsonParser{},
	".yml":  yamlParser{},
	".yaml": yamlParser{},
}

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
// timestamps to a string in the shortest form that keeps every part of the value
// (date only for midnight in UTC, date and time for other UTC values without
// fractional seconds, RFC 3339 with offset and nanoseconds otherwise). A string
// that parses as a timestamp is normalized the same way, so the same moment
// compares equal however it is spelled and whichever format it comes from.
// float64 is also the limit of comparison: integers above 2^53 that differ
// only in their low digits become the same value.
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

	root, err := extractRootMap(parsed)
	if err != nil {
		return nil, fmt.Errorf("file '%s': %w", path, err)
	}

	normalizeMap(root)

	return root, nil
}

func extractRootMap(root any) (map[string]any, error) {
	if root == nil {
		return nil, fmt.Errorf("%w: got null", ErrInvalidRoot)
	}

	m, ok := root.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%w: got %s", ErrInvalidRoot, describeRootKind(root))
	}

	return m, nil
}

func describeRootKind(root any) string {
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
	ext, err := extractExt(filePath)
	if err != nil {
		return nil, err
	}

	p, ok := extParsers[ext]
	if !ok {
		return nil, fmt.Errorf("%w: '%s' for file '%s'", ErrUnsupportedExtension, ext, filePath)
	}

	return p, nil
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
