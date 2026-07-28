package code

import (
	"code/internal/diff"
	"code/internal/formatters"
	"code/internal/parsers"
)

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
