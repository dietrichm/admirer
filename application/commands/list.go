package commands

import (
	"fmt"
	"io"

	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/services"
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
		return list(services.AvailableServices, command.OutOrStdout(), args)
	},
}

func list(serviceLoader domain.ServiceLoader, writer io.Writer, args []string) error {
	serviceName := args[0]

	service, err := serviceLoader.ForName(serviceName)
	if err != nil {
		return err
	}

	defer service.Close()

	tracks, err := service.GetLovedTracks()
	if err != nil {
		return err
	}

	for _, track := range tracks {
		fmt.Fprintln(writer, track.Artist, "-", track.Name)
	}

	return nil
}
