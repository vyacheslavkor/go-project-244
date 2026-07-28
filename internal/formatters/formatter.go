package formatters

import "code/internal/diff"

type Formatter interface {
	Format(d *diff.Diff) string
}
