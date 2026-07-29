package parsers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type parser interface {
	parse(content []byte) (map[string]any, error)
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
	// ErrPathHasNoExtension is returned when a path has no file extension.
	ErrPathHasNoExtension = errors.New("file path has no extension")
	// ErrDifferentFileFormats is returned when the two files use incompatible formats.
	ErrDifferentFileFormats = errors.New("files have different formats")
	// ErrUnsupportedExtension is returned when the file extension is not supported.
	ErrUnsupportedExtension = errors.New("unsupported extension")
	// ErrNotRegularFile is returned when the path is not a regular file.
	ErrNotRegularFile = errors.New("path is not a regular file")
	// ErrFileIsEmpty is returned when the file has zero size.
	ErrFileIsEmpty = errors.New("file is empty")
	// ErrInvalidRoot is returned when the root value is not a JSON object or YAML mapping.
	ErrInvalidRoot = errors.New("root value must be a JSON object or a YAML mapping")
	// ErrFailedToParseFile is returned when file contents cannot be decoded.
	ErrFailedToParseFile = errors.New("failed to parse file")
	// ErrFailedToReadFile is returned when the file cannot be read from disk.
	ErrFailedToReadFile = errors.New("failed to read file")
)

// ParseFiles reads and parses two configuration files and returns their
// contents as nested maps.
//
// f1 and f2 must be paths to existing non-empty regular files whose root
// value is a JSON object or YAML mapping. Supported formats are JSON
// (.json) and YAML (.yml, .yaml). Both files must use compatible formats:
// JSON with JSON, or YAML with YAML (.yml and .yaml may be mixed).
//
// Paths are validated as regular non-empty files before extension checks,
// so directories and missing paths are reported as file errors rather than
// "no extension" usage errors.
//
// On failure, err may wrap one of the package sentinel errors
// (for example [ErrUnsupportedExtension], [ErrDifferentFileFormats],
// [ErrInvalidRoot], [ErrFailedToParseFile]).
func ParseFiles(f1, f2 string) (parsed1, parsed2 map[string]any, err error) {
	if err := validateInputFile(f1); err != nil {
		return nil, nil, err
	}

	if err := validateInputFile(f2); err != nil {
		return nil, nil, err
	}

	p, err := resolveParser(f1, f2)
	if err != nil {
		return nil, nil, err
	}

	content1, err := readFileContent(f1)
	if err != nil {
		return nil, nil, err
	}

	content2, err := readFileContent(f2)
	if err != nil {
		return nil, nil, err
	}

	parsed1, err = parseFile(p, f1, content1)
	if err != nil {
		return nil, nil, err
	}

	parsed2, err = parseFile(p, f2, content2)
	if err != nil {
		return nil, nil, err
	}

	return parsed1, parsed2, nil
}

func parseFile(p parser, path string, content []byte) (map[string]any, error) {
	parsed, err := p.parse(content)
	if err == nil {
		return parsed, nil
	}

	if errors.Is(err, ErrInvalidRoot) {
		return nil, fmt.Errorf("file '%s': %w", path, err)
	}

	return nil, fmt.Errorf("%w: '%s': %w", ErrFailedToParseFile, path, err)
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

func resolveParser(file1, file2 string) (parser, error) {
	format1, err := fileFormat(file1)
	if err != nil {
		return nil, err
	}

	format2, err := fileFormat(file2)
	if err != nil {
		return nil, err
	}

	if format1 != format2 {
		return nil, fmt.Errorf("%w: files '%s' and '%s'", ErrDifferentFileFormats, file1, file2)
	}

	return formatParsers[format1], nil
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
		return "", fmt.Errorf("%w: '%s'", ErrPathHasNoExtension, path)
	}

	return ext, nil
}

func validateInputFile(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%w: '%s': %w", ErrFailedToReadFile, path, err)
	}

	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("%w: '%s'", ErrNotRegularFile, path)
	}

	if fileInfo.Size() == 0 {
		return fmt.Errorf("%w: '%s'", ErrFileIsEmpty, path)
	}

	return nil
}

func readFileContent(path string) ([]byte, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("%w: '%s': %w", ErrFailedToReadFile, path, err)
	}

	return content, nil
}
