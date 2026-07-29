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
// filepath1 and filepath2 must be paths to existing non-empty regular files
// whose root value is a JSON object or YAML mapping (arrays/sequences at the
// root are rejected). Supported input formats are JSON (.json) and YAML
// (.yml, .yaml). Both files must use compatible formats: JSON with JSON, or
// YAML with YAML (.yml and .yaml may be mixed). The parser is chosen by file
// extension; content is not sniffed.
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
// Node status values are added, removed, updated, nested, and unchanged.
// The envelope status "root" is part of the JSON wire format only; it is not
// a difference-tree node status. Library failures wrap the sentinel errors
// re-exported from this package (for example [ErrInvalidFormat],
// [ErrInvalidRoot], [ErrFileIsEmpty], [ErrFailedToParseFile]).
func GenDiff(filepath1, filepath2, format string) (string, error) {
	formatter, err := formatters.NewFormatter(format)
	if err != nil {
		return "", err
	}

	parsed1, parsed2, err := parsers.ParseFiles(filepath1, filepath2)
	if err != nil {
		return "", err
	}

	t := diff.NewTree(parsed1, parsed2)
	formatted, err := formatter.Format(t)
	if err != nil {
		return "", err
	}

	return formatted, nil
}
