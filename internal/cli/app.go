package cli

import (
	"code"
	"code/internal/formatters"
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

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

			diff, err := code.GenDiff(file1, file2, cmd.String("format"))
			if err != nil {
				return returnError(err, cmd)
			}

			fmt.Fprintln(cmd.Writer, diff)

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
				Usage:       "output format",
				DefaultText: "\"stylish\"",
				Value:       "stylish",
			},
		},
	}
}

func newArgumentsCountError(cmd *cli.Command, expected, got int) error {
	cli.ShowAppHelp(cmd)

	return cli.Exit(
		fmt.Sprintf("incorrect usage: expected %d argument, got %d", expected, got),
		1,
	)
}

func returnError(err error, cmd *cli.Command) error {
	if errors.Is(err, formatters.ErrInvalidFormat) {
		cli.ShowAppHelp(cmd)
		return cli.Exit(err.Error(), 1)
	}

	return err
}
