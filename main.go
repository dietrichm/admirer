package main

import (
	"flag"
	"os"

	"github.com/dietrichm/admirer/services/lastfm"
	"github.com/dietrichm/admirer/services/spotify"
)

func main() {
	var spotifyLogin bool
	flag.BoolVar(&spotifyLogin, "spotify", false, "Authenticate with Spotify")

	var lastfmLogin bool
	flag.BoolVar(&lastfmLogin, "lastfm", false, "Authenticate with Last.fm")

	var oauthCode string
	flag.StringVar(&oauthCode, "oauth-code", "", "OAuth code")

	flag.Parse()

	if spotifyLogin {
		spotify.Login(oauthCode)
		os.Exit(0)
	}

	if lastfmLogin {
		lastfm.Login(oauthCode)
		os.Exit(0)
	}

	flag.Usage()
	os.Exit(1)
}
