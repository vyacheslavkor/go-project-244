package formatters

import (
	"code/internal/diff"
	"fmt"
)

type Formatter interface {
	Format(d *diff.Diff) string
}

func NewFormatter(format string) (Formatter, error) {
	switch format {
	case "stylish":
		return &StylishFormatter{}, nil
	case "plain":
		return &PlainFormatter{}, nil
	}

	return nil, fmt.Errorf("invalid format: %s", format)
}
