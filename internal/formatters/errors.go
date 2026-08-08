package formatters

import "errors"

// ErrInvalidFormat is returned when the requested output format is unknown.
var ErrInvalidFormat = errors.New("invalid format")
