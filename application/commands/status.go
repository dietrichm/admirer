package commands

import (
	"fmt"
	"io"

	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/services"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(statusCommand)
}

var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "Retrieve status for services",
	RunE: func(command *cobra.Command, args []string) error {
		return status(services.AvailableServices, command.OutOrStdout())
	},
}

func status(serviceLoader domain.ServiceLoader, writer io.Writer) error {
	for _, serviceName := range serviceLoader.Names() {
		service, err := serviceLoader.ForName(serviceName)
		if err != nil {
			return err
		}

		statusForService(service, writer)
	}

	return nil
}

func statusForService(service domain.Service, writer io.Writer) {
	fmt.Fprintln(writer, service.Name())

	if !service.Authenticated() {
		fmt.Fprintln(writer, "\tNot logged in")
		return
	}

	username, err := service.GetUsername()

	if err != nil {
		fmt.Fprintln(writer, "\tNot authenticated -", err.Error())
	} else {
		fmt.Fprintln(writer, "\tAuthenticated as", username)
	}
}
