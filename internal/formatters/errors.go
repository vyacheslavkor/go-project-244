package formatters

import "errors"

var (
	// ErrInvalidFormat is returned when the requested output format is unknown.
	ErrInvalidFormat = errors.New("invalid format")
	// ErrFormatOutput is returned when the difference tree cannot be rendered as output.
	ErrFormatOutput = errors.New("failed to format output")
)
