package formatters

import (
	"code/internal/diff"
	"fmt"
)

type Formatter interface {
	Format(d *diff.Diff) (string, error)
}

func NewFormatter(format string) (Formatter, error) {
	switch format {
	case "stylish":
		return stylishFormatter{}, nil
	case "plain":
		return plainFormatter{}, nil
	case "json":
		return jsonFormatter{}, nil
	}

	return nil, fmt.Errorf("invalid format: %s", format)
}
