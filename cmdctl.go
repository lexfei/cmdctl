package main

import (
	"os"

	"cmdctl/cmd"
	cmdutil "cmdctl/cmd/util"
)

func main() {
	cmd := cmd.NewCommand(cmdutil.NewFactory(), os.Stdin, os.Stdout, os.Stderr)
	if cmd.Execute() != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
