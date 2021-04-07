package commands

import (
	"io"

	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/services"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(syncCommand)
}

var syncCommand = &cobra.Command{
	Use:   "sync <source-service> <target-service>",
	Short: "Sync recently loved tracks from one service to another",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(command *cobra.Command, args []string) error {
		return sync(services.AvailableServices, command.OutOrStdout(), args)
	},
}

func sync(serviceLoader domain.ServiceLoader, writer io.Writer, args []string) error {
	return nil
}
