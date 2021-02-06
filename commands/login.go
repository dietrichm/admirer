package commands

import (
	"fmt"

	"github.com/dietrichm/admirer/services"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(loginCommand)
}

var loginCommand = &cobra.Command{
	Use:   "login <service> [oauth-code]",
	Short: "Log in on external service",
	Args:  cobra.MinimumNArgs(1),
	Run: func(command *cobra.Command, args []string) {
		service, err := services.ForName(args[0])
		if err != nil {
			panic(err)
		}

		if len(args) < 2 {
			fmt.Println(service.Name() + " authentication URL: " + service.CreateAuthURL())
			return
		}

		if err := service.Authenticate(args[1]); err != nil {
			panic(err)
		}

		username, err := service.GetUsername()
		if err != nil {
			panic(err)
		}

		fmt.Println("Logged in on " + service.Name() + " as " + username)
	},
}
