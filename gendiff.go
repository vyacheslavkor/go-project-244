package code

import (
	"code/internal/diff"
	"code/internal/formatters"
	"code/internal/parsers"
)

func GenDiff(filepath1, filepath2, format string) (string, error) {
	if err := formatters.ValidateFormat(format); err != nil {
		return "", err
	}

	parsed1, parsed2, err := parsers.ParseFiles(filepath1, filepath2)
	if err != nil {
		return "", err
	}

	d := diff.NewDiff(parsed1, parsed2)
	formatted := formatters.ToStylish(d)

	return formatted, nil
}
