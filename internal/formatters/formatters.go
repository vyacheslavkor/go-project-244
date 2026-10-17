package formatters

import (
	"code/internal/diff"
	"fmt"
	"strconv"
)

// Format names.
const (
	FormatStylish = "stylish"
	FormatPlain   = "plain"
	FormatJSON    = "json"
)

// Formatter formats a difference tree as a string.
type Formatter interface {
	// Format renders root according to the formatter implementation.
	// It panics if a node carries a status the implementation does not
	// handle, which cannot happen for a tree built by [diff.Build].
	Format(root *diff.Node) (string, error)
}

// NewFormatter returns a Formatter for the given format name
// ([FormatStylish], [FormatPlain], or [FormatJSON]).
// An empty format is treated as [FormatStylish].
//
// Empty document rendering (both roots have no properties):
//   - stylish: pretty-printed empty object "{}"
//   - plain: empty string
//   - json: compact single-line root object with "children" omitted
//
// Identical files are a different case, and not a special one: every property
// is unchanged, so stylish renders the whole tree, json renders the full
// "children" list of unchanged nodes, and plain renders an empty string.
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

func formatNumber(val float64) string {
	return strconv.FormatFloat(val, 'f', -1, 64)
}
