package main

import (
	"flag"
	"os"

	"github.com/dietrichm/admirer/services/spotify"
)

func main() {
	var spotifyLogin bool
	flag.BoolVar(&spotifyLogin, "spotify", false, "Authenticate with Spotify")

	var oauthCode string
	flag.StringVar(&oauthCode, "oauth-code", "", "OAuth code")

	flag.Parse()

	if spotifyLogin {
		spotify.Login(oauthCode)
		os.Exit(0)
	}

	flag.Usage()
	os.Exit(1)
}
