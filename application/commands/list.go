package commands

import (
	"fmt"
	"io"

	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/services"
	"github.com/spf13/cobra"
)

func init() {
	listCommand.Flags().IntVarP(&limit, "limit", "l", 10, "Limit number of returned tracks")
	rootCommand.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:   "list <service>",
	Short: "List loved tracks on specified service",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(command *cobra.Command, args []string) error {
		return list(services.AvailableServices, limit, command.OutOrStdout(), args)
	},
}

func list(serviceLoader domain.ServiceLoader, limit int, writer io.Writer, args []string) error {
	serviceName := args[0]

	service, err := serviceLoader.ForName(serviceName)
	if err != nil {
		return err
	}

	defer service.Close()

	if !service.Authenticated() {
		return fmt.Errorf("not logged in on %s", service.Name())
	}

	tracks, err := service.GetLovedTracks(limit)
	if err != nil {
		return err
	}

	for _, track := range tracks {
		fmt.Fprintln(writer, track.String())
	}

	return nil
}
