package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:   "list <service>",
	Short: "List loved tracks on specified service",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(command *cobra.Command, args []string) error {
		return nil
	},
}
