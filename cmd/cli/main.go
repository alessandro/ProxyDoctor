package main

import (
	"os"

	"github.com/francomano/proxydoctor/cmd/cli/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
