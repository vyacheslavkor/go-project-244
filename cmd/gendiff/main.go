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
		fmt.Fprintln(command.ErrWriter, err)
		os.Exit(1)
	}
}
