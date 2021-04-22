package commands

import (
	"fmt"
	"io"

	"github.com/dietrichm/admirer/domain"
	"github.com/dietrichm/admirer/infrastructure/authentication"
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
		return login(services.AvailableServices, authentication.DefaultCallbackProvider, command.OutOrStdout(), args)
	},
}

func login(serviceLoader domain.ServiceLoader, callbackProvider authentication.CallbackProvider, writer io.Writer, args []string) error {
	serviceName := args[0]

	service, err := serviceLoader.ForName(serviceName)
	if err != nil {
		return err
	}

	defer service.Close()
	redirectURL := "https://admirer.test"

	if len(args) < 2 {
		fmt.Fprintln(writer, service.Name(), "authentication URL:", service.CreateAuthURL(redirectURL))
		return nil
	}

	if err := service.Authenticate(args[1], redirectURL); err != nil {
		return err
	}

	username, err := service.GetUsername()
	if err != nil {
		return err
	}

	fmt.Fprintln(writer, "Logged in on", service.Name(), "as", username)
	return nil
}
