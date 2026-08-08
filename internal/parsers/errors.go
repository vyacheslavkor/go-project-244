package parsers

import "errors"

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
