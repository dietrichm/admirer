package commands

import (
	"fmt"
	"io"

	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/services"
	"github.com/spf13/cobra"
)

func init() {
	syncCommand.Flags().IntVarP(&limit, "limit", "l", 10, "Limit number of synced tracks")
	rootCommand.AddCommand(syncCommand)
}

var syncCommand = &cobra.Command{
	Use:   "sync <source-service> <target-service>",
	Short: "Sync recently loved tracks from one service to another",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(command *cobra.Command, args []string) error {
		return sync(services.AvailableServices, limit, command.OutOrStdout(), args)
	},
}

func sync(serviceLoader domain.ServiceLoader, limit int, writer io.Writer, args []string) error {
	sourceServiceName := args[0]
	targetServiceName := args[1]

	sourceService, err := serviceLoader.ForName(sourceServiceName)
	if err != nil {
		return err
	}

	targetService, err := serviceLoader.ForName(targetServiceName)
	if err != nil {
		return err
	}

	defer sourceService.Close()
	defer targetService.Close()

	if !sourceService.Authenticated() {
		return fmt.Errorf("not logged in on %q", sourceServiceName)
	}

	if !targetService.Authenticated() {
		return fmt.Errorf("not logged in on %q", targetServiceName)
	}

	tracks, err := sourceService.GetLovedTracks(limit)
	if err != nil {
		return err
	}

	for _, track := range tracks {
		if err := targetService.LoveTrack(track); err != nil {
			return err
		}

		fmt.Fprintln(writer, "Synced:", track.String())
	}

	return nil
}
