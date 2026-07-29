package formatters

import (
	"code/internal/diff"
	"errors"
	"fmt"
)

const (
	FormatStylish = "stylish"
	FormatPlain   = "plain"
	FormatJSON    = "json"
)

// Formatter formats a difference tree as a string.
type Formatter interface {
	// Format renders t according to the formatter implementation.
	Format(t *diff.Tree) (string, error)
}

// ErrInvalidFormat is returned when the requested output format is unknown.
var ErrInvalidFormat = errors.New("invalid format")

// NewFormatter returns a Formatter for the given format name
// ([FormatStylish], [FormatPlain], or [FormatJSON]).
// An empty format is treated as [FormatStylish].
//
// No-diff rendering:
//   - stylish: pretty-printed empty object "{}"
//   - plain: empty string
//   - json: compact single-line root object; "children" omitted when empty;
//     unchanged nodes are included when present
func NewFormatter(format string) (Formatter, error) {
	if format == "" {
		format = FormatStylish
	}

	switch format {
	case FormatStylish:
		return stylishFormatter{}, nil
	case FormatPlain:
		return plainFormatter{}, nil
	case FormatJSON:
		return jsonFormatter{}, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrInvalidFormat, format)
}
