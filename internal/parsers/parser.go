package parsers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type formatParser interface {
	parse(content []byte) (parsed map[string]any, err error)
}

var extHandlers = map[string]formatParser{
	".json": &jsonParser{},
	".yml":  &yamlParser{},
	".yaml": &yamlParser{},
}

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
	parser1, err := resolveHandler(file1)
	if err != nil {
		return nil, err
	}

	parser2, err := resolveHandler(file2)
	if err != nil {
		return nil, err
	}

	if parser1 != parser2 {
		return nil, fmt.Errorf("files '%s' and '%s' have different formats", file1, file2)
	}

	return parser1, nil
}

func resolveHandler(file string) (formatParser, error) {
	ext := strings.ToLower(filepath.Ext(file))
	if ext == "" {
		return nil, fmt.Errorf("file '%s' has no extension", file)
	}

	handler, ok := extHandlers[ext]
	if !ok {
		return nil, fmt.Errorf("file '%s' has unsupported extension '%s'", file, ext)
	}

	return handler, nil
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
