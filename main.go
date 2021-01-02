package main

import (
	"flag"
	"os"

	"github.com/dietrichm/admirer/commands"
)

func main() {
	var loginServiceName string
	flag.StringVar(&loginServiceName, "login", "", "Log in on external service")

	var oauthCode string
	flag.StringVar(&oauthCode, "oauth-code", "", "OAuth code")

	flag.Parse()

	if len(loginServiceName) != 0 {
		commands.Login(loginServiceName, oauthCode)
		os.Exit(0)
	}

	flag.Usage()
	os.Exit(1)
}
