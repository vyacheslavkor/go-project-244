package code

import (
	"code/internal/diff"
	"code/internal/formatters"
	"code/internal/parsers"
)

// GenDiff compares two configuration files and returns their difference
// as a formatted string.
//
// filepath1 and filepath2 must be paths to existing non-empty regular files.
// Supported input formats are JSON (.json) and YAML (.yml, .yaml).
// Both files must use compatible formats: JSON with JSON, or YAML with YAML
// (.yml and .yaml may be mixed).
//
// format selects the output style: "stylish" (default-style tree), "plain",
// or "json".
func GenDiff(filepath1, filepath2, format string) (string, error) {
	formatter, err := formatters.NewFormatter(format)
	if err != nil {
		return "", err
	}

	parsed1, parsed2, err := parsers.ParseFiles(filepath1, filepath2)
	if err != nil {
		return "", err
	}

	d := diff.NewTree(parsed1, parsed2)
	formatted, err := formatter.Format(d)
	if err != nil {
		return "", err
	}

	return formatted, nil
}
