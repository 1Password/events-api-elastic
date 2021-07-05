package main

import (
	"os"

	"go.1password.io/eventsapibeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}