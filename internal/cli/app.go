package cli

import (
	"code"
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

// NewCommand builds the gendiff CLI command with flags and actions.
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		ArgsUsage: "<filepath1> <filepath2>",
		Description: `Compare two configuration files and print their difference.

Arguments:
  <filepath1>, <filepath2>  Paths to existing non-empty regular files. Exactly two
                positional arguments are required.

Input constraints:
  The root value of each file must be a JSON object or a YAML mapping.
  Root-level arrays/sequences, scalars, and null are rejected.
  Empty files are rejected.

Input formats:
  JSON  (.json)
  YAML  (.yml, .yaml)

Both files must use compatible formats: JSON with JSON, or YAML with YAML.
Mixing .yml and .yaml is allowed; mixing JSON and YAML is not.
The parser is chosen by file extension (content is not sniffed).

Output formats (--format / -f):
  stylish  Nested tree with +/- markers (default).
           No changes: pretty-printed empty object.
  plain    Line-oriented human-readable messages (skips unchanged).
           No changes: empty output (nothing is printed).
  json     Compact single-line JSON tree of changes.
           Envelope: {"key":"","status":"root","children":[...]}.
           Field "children" is omitted when there are no nodes.
           Node fields: key, status, old_value?, value?, children?.
           Node statuses: root (document root, empty key),
           added, removed, updated, nested, unchanged.
           Unchanged nodes are included.

Usage errors (wrong args/flags/input contract) print this help on stdout
and a reason on stderr. Operational errors (missing/unreadable files,
malformed content) print only the reason on stderr.

Examples:
  gendiff file1.json file2.json
  gendiff -f plain a.yml b.yaml
  gendiff --format=json before.json after.json`,
		Action: func(_ context.Context, cmd *cli.Command) error {
			const expectedArgsCount = 2
			if cmd.Args().Len() != expectedArgsCount {
				return newArgumentsCountError(cmd, expectedArgsCount, cmd.Args().Len())
			}

			beforePath, afterPath := cmd.Args().Get(0), cmd.Args().Get(1)
			out, err := code.GenDiff(beforePath, afterPath, cmd.String("format"))
			if err != nil {
				return returnError(err, cmd)
			}

			if out != "" {
				_, err := fmt.Fprintln(cmd.Writer, out)
				if err != nil {
					return fmt.Errorf("failed to write output: %w", err)
				}
			}

			return nil
		},
		OnUsageError: func(_ context.Context, cmd *cli.Command, err error, _ bool) error {
			cli.ShowAppHelp(cmd)

			return cli.Exit(fmt.Sprintf("incorrect usage: %v", err), 1)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "format",
				Aliases:     []string{"f"},
				Usage:       "output format (allowed: stylish, plain, json)",
				DefaultText: "\"stylish\"",
				Value:       code.FormatStylish,
			},
		},
	}
}

func newArgumentsCountError(cmd *cli.Command, expected, got int) error {
	cli.ShowAppHelp(cmd)

	return cli.Exit(
		fmt.Sprintf("incorrect usage: expected %d arguments, got %d", expected, got),
		1,
	)
}

func returnError(err error, cmd *cli.Command) error {
	if !isUsageError(err) {
		return err
	}

	cli.ShowAppHelp(cmd)

	return cli.Exit(fmt.Sprintf("incorrect usage: %v", err), 1)
}

func isUsageError(err error) bool {
	return errors.Is(err, code.ErrInvalidFormat) ||
		errors.Is(err, code.ErrMissingExtension) ||
		errors.Is(err, code.ErrUnsupportedExtension) ||
		errors.Is(err, code.ErrFormatMismatch) ||
		errors.Is(err, code.ErrEmptyFile) ||
		errors.Is(err, code.ErrNotRegularFile) ||
		errors.Is(err, code.ErrInvalidRoot)
}
