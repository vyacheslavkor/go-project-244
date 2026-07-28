package formatters

import "fmt"

func ValidateFormat(format string) error {
	if format == "stylish" {
		return nil
	}
	return fmt.Errorf("invalid format: %s", format)
}
