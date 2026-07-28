package parsers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type formatParser interface {
	parse(content []byte) (parsed map[string]any, err error)
}

var extHandlers = map[string]formatParser{
	".json": jsonParser{},
	".yml":  yamlParser{},
	".yaml": yamlParser{},
}

var (
	ErrPathHasNoExtension   = errors.New("cannot extract extension from path")
	ErrDifferentFileFormats = errors.New("files have different formats")
	ErrUnsupportedExtension = errors.New("unsupported extension")
)

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
		return nil, nil, fmt.Errorf("failed to parse file '%s': %w", f1, err)
	}

	parsed2, err = p.parse(content2)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse file '%s': %w", f2, err)
	}

	return parsed1, parsed2, nil
}

func resolveParser(file1, file2 string) (formatParser, error) {
	ext1, err := extractExt(file1)
	if err != nil {
		return nil, err
	}

	if !isAllowedExt(ext1) {
		return nil, fmt.Errorf("%w: '%s' for file '%s'", ErrUnsupportedExtension, ext1, file1)
	}

	ext2, err := extractExt(file2)
	if err != nil {
		return nil, err
	}

	if !isAllowedExt(ext2) {
		return nil, fmt.Errorf("%w: '%s' for file '%s'", ErrUnsupportedExtension, ext2, file2)
	}

	if !isCompatibleExt(ext1, ext2) {
		return nil, fmt.Errorf("%w: files '%s' and '%s'", ErrDifferentFileFormats, file1, file2)
	}

	return extHandlers[ext1], nil
}

func extractExt(file string) (string, error) {
	ext := strings.ToLower(filepath.Ext(file))
	if ext == "" {
		return "", fmt.Errorf("%w: %s", ErrPathHasNoExtension, file)
	}

	return ext, nil
}

func isAllowedExt(ext string) bool {
	_, ok := extHandlers[ext]
	return ok
}

func isCompatibleExt(ext1, ext2 string) bool {
	compatibleMap := map[string][]string{
		".json": {".json"},
		".yml":  {".yml", ".yaml"},
		".yaml": {".yaml", ".yml"},
	}

	compatible := compatibleMap[ext1]

	return slices.Contains(compatible, ext2)
}

func getFileContent(f string) ([]byte, error) {
	fileInfo, err := os.Stat(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read path metadata for %s: %w", f, err)
	}

	mode := fileInfo.Mode()
	if !mode.IsRegular() {
		return nil, fmt.Errorf("file '%s' is not a regular file", f)
	}

	if fileInfo.Size() == 0 {
		return nil, fmt.Errorf("file '%s' is empty", f)
	}

	content, err := os.ReadFile(filepath.Clean(f))
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s': %w", f, err)
	}

	return content, nil
}
