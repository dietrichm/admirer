package commands

import (
	"fmt"
	"io"

	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/services"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(loginCommand)
}

var loginCommand = &cobra.Command{
	Use:   "login <service> [oauth-code]",
	Short: "Log in on external service",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(command *cobra.Command, args []string) error {
		return login(services.AvailableServices, command.OutOrStdout(), args)
	},
}

func login(serviceLoader domain.ServiceLoader, writer io.Writer, args []string) error {
	serviceName := args[0]

	service, err := serviceLoader.ForName(serviceName)
	if err != nil {
		return err
	}

	if len(args) < 2 {
		fmt.Fprintln(writer, service.Name(), "authentication URL:", service.CreateAuthURL())
		return nil
	}

	if err := service.Authenticate(args[1]); err != nil {
		return err
	}

	username, err := service.GetUsername()
	if err != nil {
		return err
	}

	fmt.Fprintln(writer, "Logged in on", service.Name(), "as", username)
	return nil
}
