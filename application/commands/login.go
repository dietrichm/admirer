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
		return login(services.AvailableServices, authentication.DefaultCallbackServer, command.OutOrStdout(), args)
	},
}

func login(serviceLoader domain.ServiceLoader, callbackServer authentication.CallbackServer, writer io.Writer, args []string) error {
	serviceName := args[0]

	service, err := serviceLoader.ForName(serviceName)
	if err != nil {
		return err
	}

	defer service.Close()
	redirectURL := "http://0.0.0.0:8080/"

	if len(args) < 2 {
		fmt.Fprintln(writer, service.Name(), "authentication URL:", service.CreateAuthURL(redirectURL))
		fmt.Fprintln(writer, "Waiting for authentication callback...")
	}

	token, err := callbackServer.ReadCode(service.CodeParam())
	if err != nil {
		return fmt.Errorf("failed to read authentication code: %w", err)
	}

	if err := service.Authenticate(token, redirectURL); err != nil {
		return err
	}

	username, err := service.GetUsername()
	if err != nil {
		return err
	}

	fmt.Fprintln(writer, "Logged in on", service.Name(), "as", username)
	return nil
}
