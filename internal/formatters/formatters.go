package formatters

import (
	"code/internal/diff"
	"fmt"
)

type Formatter interface {
	Format(d *diff.Diff) string
}

func NewFormatter(format string) (Formatter, error) {
	if format == "stylish" {
		return &StylishFormatter{}, nil
	}

	return nil, fmt.Errorf("invalid format: %s", format)
}
