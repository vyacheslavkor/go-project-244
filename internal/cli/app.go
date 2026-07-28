package cli

import (
	"code"
	"code/internal/formatters"
	"code/internal/parsers"
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

// NewCommand builds the gendiff CLI command with flags and actions.
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "gendiff",
		Usage: "Compares two configuration files and shows a difference.",
		Action: func(_ context.Context, cmd *cli.Command) error {
			const expectedArgsCount = 2
			if cmd.Args().Len() != expectedArgsCount {
				return newArgumentsCountError(cmd, expectedArgsCount, cmd.Args().Len())
			}

			file1, file2 := cmd.Args().Get(0), cmd.Args().Get(1)

			out, err := code.GenDiff(file1, file2, cmd.String("format"))
			if err != nil {
				return returnError(err, cmd)
			}

			fmt.Fprintln(cmd.Writer, out)

			return nil
		},
		OnUsageError: func(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error {
			cli.ShowAppHelp(cmd)

			return cli.Exit(fmt.Sprintf("incorrect usage: %v", err), 1)
		},
		ArgsUsage: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "format",
				Aliases:     []string{"f"},
				Usage:       "output format (allowed: stylish, plain, json)",
				DefaultText: "\"stylish\"",
				Value:       "stylish",
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
	return errors.Is(err, formatters.ErrInvalidFormat) ||
		errors.Is(err, parsers.ErrPathHasNoExtension) ||
		errors.Is(err, parsers.ErrUnsupportedExtension) ||
		errors.Is(err, parsers.ErrDifferentFileFormats)
}
