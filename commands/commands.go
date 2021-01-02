package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "admirer",
	Short: "A command line utility to sync song likes between Spotify and Last.fm.",
}

// Execute runs the requested CLI command.
func Execute() {
	err := rootCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}
