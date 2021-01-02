package commands

import (
	"fmt"

	"github.com/dietrichm/admirer/services"
	"github.com/spf13/cobra"
)

// Login runs the command for logging in on an external service.
func Login(serviceName string, oauthCode string) {
	service := services.ForName(serviceName)

	if len(oauthCode) == 0 {
		fmt.Println(service.Name() + " authentication URL: " + service.CreateAuthURL())
		return
	}

	service.Authenticate(oauthCode)

	fmt.Println("Logged in on " + service.Name() + " as " + service.GetUsername())
}

var rootCommand = &cobra.Command{
	Use:   "admirer",
	Short: "A command line utility to sync song likes between Spotify and Last.fm.",
}

// Execute runs the requested CLI command.
func Execute() {
	err := rootCommand.Execute()
	if err != nil {
		panic(err)
	}
}
