package code

import (
	"code/internal/formatters"
	"code/internal/parsers"
)

var (
	// ErrInvalidFormat is returned when the output format name is unknown.
	ErrInvalidFormat = formatters.ErrInvalidFormat
	// ErrPathHasNoExtension is returned when a path has no file extension.
	ErrPathHasNoExtension = parsers.ErrPathHasNoExtension
	// ErrDifferentFileFormats is returned when the two files use incompatible formats.
	ErrDifferentFileFormats = parsers.ErrDifferentFileFormats
	// ErrUnsupportedExtension is returned when the file extension is not supported.
	ErrUnsupportedExtension = parsers.ErrUnsupportedExtension
	// ErrNotRegularFile is returned when the path is not a regular file.
	ErrNotRegularFile = parsers.ErrNotRegularFile
	// ErrFileIsEmpty is returned when the file has zero size.
	ErrFileIsEmpty = parsers.ErrFileIsEmpty
	// ErrFailedToParseFile is returned when file contents cannot be decoded.
	ErrFailedToParseFile = parsers.ErrFailedToParseFile
	// ErrFailedToReadFile is returned when the file cannot be read from disk.
	ErrFailedToReadFile = parsers.ErrFailedToReadFile
)
