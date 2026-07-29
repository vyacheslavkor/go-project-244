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
	ErrPathHasNoExtension = errors.New("cannot extract extension from path")
	// ErrDifferentFileFormats is returned when the two files use incompatible formats.
	ErrDifferentFileFormats = errors.New("files have different formats")
	// ErrUnsupportedExtension is returned when the file extension is not supported.
	ErrUnsupportedExtension = errors.New("unsupported extension")
	// ErrNotRegularFile is returned when the path is not a regular file.
	ErrNotRegularFile = errors.New("file is not a regular file")
	// ErrFileIsEmpty is returned when the file has zero size.
	ErrFileIsEmpty = errors.New("file is empty")
	// ErrFailedToParseFile is returned when file contents cannot be decoded.
	ErrFailedToParseFile = errors.New("failed to parse file")
	// ErrFailedToReadFile is returned when the file cannot be read from disk.
	ErrFailedToReadFile = errors.New("failed to read file")
)

// ParseFiles reads and parses two configuration files and returns their
// contents as nested maps.
//
// f1 and f2 must be paths to existing non-empty regular files.
// Supported formats are JSON (.json) and YAML (.yml, .yaml).
// Both files must use compatible formats: JSON with JSON, or YAML with YAML
// (.yml and .yaml may be mixed).
//
// On failure, err may wrap one of the package sentinel errors
// (for example [ErrUnsupportedExtension], [ErrDifferentFileFormats],
// [ErrFailedToParseFile]).
func ParseFiles(f1, f2 string) (parsed1, parsed2 map[string]any, err error) {
	p, err := resolveParser(f1, f2)
	if err != nil {
		return nil, nil, err
	}

	content1, err := getFileContent(f1)
	if err != nil {
		return nil, nil, err
	}

	content2, err := getFileContent(f2)
	if err != nil {
		return nil, nil, err
	}

	parsed1, err = p.parse(content1)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: '%s': %w", ErrFailedToParseFile, f1, err)
	}

	parsed2, err = p.parse(content2)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: '%s': %w", ErrFailedToParseFile, f2, err)
	}

	return parsed1, parsed2, nil
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
		return "", fmt.Errorf("%w: %s", ErrPathHasNoExtension, path)
	}

	return ext, nil
}

func getFileContent(path string) ([]byte, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read path metadata for %s: %w", path, err)
	}

	mode := fileInfo.Mode()
	if !mode.IsRegular() {
		return nil, fmt.Errorf("%w: '%s'", ErrNotRegularFile, path)
	}

	if fileInfo.Size() == 0 {
		return nil, fmt.Errorf("%w: '%s'", ErrFileIsEmpty, path)
	}

	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("%w: '%s': %w", ErrFailedToReadFile, path, err)
	}

	return content, nil
}
