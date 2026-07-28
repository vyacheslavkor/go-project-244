package formatters

import (
	"code/internal/diff"
	"errors"
	"fmt"
)

// Formatter formats a difference tree as a string.
type Formatter interface {
	// Format renders t according to the formatter implementation.
	Format(t *diff.Tree) (string, error)
}

// ErrInvalidFormat is returned when the requested output format is unknown.
var ErrInvalidFormat = errors.New("invalid format")

// NewFormatter returns a Formatter for the given format name
// ("stylish", "plain", or "json").
func NewFormatter(format string) (Formatter, error) {
	switch format {
	case "stylish":
		return stylishFormatter{}, nil
	case "plain":
		return plainFormatter{}, nil
	case "json":
		return jsonFormatter{}, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrInvalidFormat, format)
}
