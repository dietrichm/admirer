package commands

import (
	"fmt"

	"github.com/dietrichm/admirer/services"
)

// Login runs the command for logging in on an external service.
func Login(serviceName string, oauthCode string) {
	service := services.ForName(serviceName)

	if len(oauthCode) == 0 {
		fmt.Println(serviceName + " authentication URL: " + service.CreateAuthURL())
		return
	}

	service.Authenticate(oauthCode)

	fmt.Println("Logged in on " + serviceName + " as " + service.GetUsername())
}
