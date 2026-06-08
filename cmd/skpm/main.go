package main

import (
	"os"

	"github.com/domehahn/skpm/internal/cli"
)

func main() {
	root := cli.NewRootCmd()
	if err := root.Execute(); err != nil {
		cli.HandleError(err, cli.OutputText)
		os.Exit(1)
	}
}
