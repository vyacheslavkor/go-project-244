package cli

import (
	"code"
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

const commandDescription = `Compare two configuration files and print their difference.

Arguments:
  <filepath1>, <filepath2>  Paths to existing non-empty regular files. Exactly two
                positional arguments are required.

Input constraints:
  The root value of each file must be a JSON object or a YAML mapping
  with string keys. Root-level arrays/sequences, scalars, and null are
  rejected, and so is a YAML mapping whose keys are not strings.
  Empty files are rejected.

Input formats:
  JSON  (.json)
  YAML  (.yml, .yaml)

The two files need not be in the same format; each one is parsed by its
own extension. The parser is chosen by extension, not by inspecting the
file.

Value comparison:
  Values are compared in a normalized form, so the same value written
  differently matches: numbers become float64 (1 equals 1.0), and a
  timestamp -- including a string that parses as one -- becomes the
  shortest string that keeps every part of the value. One moment matches
  however it is spelled and whichever format it comes from; values that
  differ as moments (fractional seconds, zone offsets) stay different.
  Comparison carries float64 precision, so integers that differ only
  beyond it (above 2^53) become the same value and can render
  differently from the file text.

Value rendering:
  Numbers are printed in positional notation in every format, so a file
  that says 1000000 shows 1000000 and never 1e+06. A value looks the
  same whichever format is selected.

Output formats (--format / -f):
  stylish  Nested tree with +/- markers (default).
  plain    Line-oriented human-readable messages (skips unchanged).
  json     Compact single-line JSON tree.
           Envelope: {"key":"","status":"root","children":[...]}.
           Field "children" is omitted when there are no nodes.
           Node fields: key, status, old_value?, value?, children?.
           Node statuses: root (document root, empty key),
           added, removed, updated, nested, unchanged.
           Unchanged nodes are included.

Identical files are not a special case: stylish prints the whole tree with
every property unchanged, json prints the full "children" list of unchanged
nodes, and plain prints nothing, because it skips unchanged properties.

A document with no properties (both roots are {}) gives an empty tree:
stylish prints a pretty-printed empty object, json prints
{"key":"","status":"root"} with "children" omitted, and plain prints
nothing. An empty file is a different case and is rejected.

Usage errors (wrong args/flags/input contract) print this help on stdout
and a reason on stderr. Operational errors (missing/unreadable files,
malformed content) print only the reason on stderr.

Examples:
  gendiff file1.json file2.json
  gendiff -f plain a.yml b.yaml
  gendiff --format=json before.json after.json`

// NewCommand builds the gendiff CLI command with flags and actions.
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:        "gendiff",
		Usage:       "Compares two configuration files and shows a difference.",
		ArgsUsage:   "<filepath1> <filepath2>",
		Description: commandDescription,
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
		errors.Is(err, code.ErrEmptyFile) ||
		errors.Is(err, code.ErrNotRegularFile) ||
		errors.Is(err, code.ErrInvalidRoot)
}
