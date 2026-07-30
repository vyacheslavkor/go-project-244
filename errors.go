package code

import (
	"code/internal/formatters"
	"code/internal/parsers"
)

var (
	// ErrInvalidFormat is returned when the output format name is unknown.
	ErrInvalidFormat = formatters.ErrInvalidFormat
	// ErrMissingExtension is returned when a path has no file extension.
	ErrMissingExtension = parsers.ErrMissingExtension
	// ErrFormatMismatch is returned when the two files use incompatible formats.
	ErrFormatMismatch = parsers.ErrFormatMismatch
	// ErrUnsupportedExtension is returned when the file extension is not supported.
	ErrUnsupportedExtension = parsers.ErrUnsupportedExtension
	// ErrNotRegularFile is returned when the path is not a regular file.
	ErrNotRegularFile = parsers.ErrNotRegularFile
	// ErrEmptyFile is returned when the file has zero size.
	ErrEmptyFile = parsers.ErrEmptyFile
	// ErrInvalidRoot is returned when the root value is not a JSON object or YAML mapping.
	ErrInvalidRoot = parsers.ErrInvalidRoot
	// ErrParseFile is returned when file contents cannot be decoded.
	ErrParseFile = parsers.ErrParseFile
	// ErrReadFile is returned when the file cannot be read from disk.
	ErrReadFile = parsers.ErrReadFile
)
