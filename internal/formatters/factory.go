package formatters

import (
	"fmt"
)

func NewFormatter(format string) (Formatter, error) {
	if format == "stylish" {
		return &StylishFormatter{}, nil
	}

	return nil, fmt.Errorf("invalid format: %s", format)
}
