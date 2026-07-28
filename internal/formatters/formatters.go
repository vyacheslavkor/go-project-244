package formatters

import (
	"code/internal/diff"
	"errors"
	"fmt"
)

type Formatter interface {
	Format(d *diff.Tree) (string, error)
}

var ErrInvalidFormat = errors.New("invalid format")

func NewFormatter(format string) (Formatter, error) {
	switch format {
	case "stylish":
		return stylishFormatter{}, nil
	case "plain":
		return plainFormatter{}, nil
	case "json":
		return jsonFormatter{}, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrInvalidFormat, format)
}
