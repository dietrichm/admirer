package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(statusCommand)
}

var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "Retrieve status for services",
	RunE: func(command *cobra.Command, args []string) error {
		return nil
	},
}
