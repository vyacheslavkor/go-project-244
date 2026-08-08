package main

import (
	"code/internal/cli"
	"context"
	"fmt"
	"os"
)

func main() {
	command := cli.NewCommand()
	if err := command.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(command.ErrWriter, err) //nolint:errcheck // diagnostic before os.Exit(1)
		os.Exit(1)
	}
}
