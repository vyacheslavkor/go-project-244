// Package code provides GenDiff, which compares two configuration files
// (JSON or YAML) and returns their difference as a formatted string.
package code

import (
	"code/internal/diff"
	"code/internal/formatters"
	"code/internal/parsers"
)

// GenDiff compares two configuration files and returns their difference
// as a formatted string.
//
// beforePath and afterPath must be paths to existing non-empty regular files
// whose root value is a JSON object or YAML mapping (arrays/sequences at the
// root are rejected). Supported input formats are JSON (.json) and YAML
// (.yml, .yaml). The files can use different formats. The parser is chosen by
// file extension, not by inspecting the file.
//
// Values are compared in a normalized form, so the same value written
// differently matches: numbers become float64 (1 equals 1.0), and a timestamp
// (including a string that parses as one) becomes the shortest string that
// keeps every part of the value. One moment therefore matches however it is
// spelled and whichever format it comes from, while values that differ as
// moments -- fractional seconds, zone offsets -- stay different. Comparison
// carries float64 precision, so integers that differ only beyond it (above
// 2^53) become the same value and can render differently from the file text.
//
// Numbers are printed in positional notation in every format, so a file that
// says 1000000 renders as 1000000 and never in exponential form. A value looks
// the same whichever format is selected.
//
// format selects the output style: [FormatStylish], [FormatPlain], or
// [FormatJSON]. An empty format is treated as [FormatStylish].
//
// When there are no differences:
//   - stylish renders a pretty-printed empty object;
//   - plain returns an empty string (no property lines; the CLI prints nothing);
//   - json emits a compact single-line root object; "children" is omitted when empty.
//
// JSON output shape (compact single-line machine-readable format):
//
//	{"key":"","status":"root","children":[...nodes]}
//
// Each node may include key, status, old_value, value, and children.
// Node status values are root (document root, empty key), added, removed,
// updated, nested, and unchanged.
// Library failures wrap the sentinel errors
// re-exported from this package (for example [ErrInvalidFormat],
// [ErrInvalidRoot], [ErrEmptyFile], [ErrParseFile]).
func GenDiff(beforePath, afterPath, format string) (string, error) {
	formatter, err := formatters.NewFormatter(format)
	if err != nil {
		return "", err
	}

	before, after, err := parsers.ParseFiles(beforePath, afterPath)
	if err != nil {
		return "", err
	}

	root := diff.Build(before, after)
	formatted, err := formatter.Format(root)
	if err != nil {
		return "", err
	}

	return formatted, nil
}
